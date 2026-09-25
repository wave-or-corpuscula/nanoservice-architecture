package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"userservice/internal/database"
	"userservice/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateOrderRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type GetUserQueryParams struct {
	WithOrders bool `form:"with_orders"`
}

func (h *Handler) GetUser(c *gin.Context) {
	var query GetUserQueryParams

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paramID := c.Param("id")
	id, err := utils.ValidateID(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	cacheKey := fmt.Sprintf("user:%d", id)

	if cached, err := h.cacher.Get(ctx, cacheKey); err == nil && cached != "" {
		var user database.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			c.JSON(http.StatusOK, &user)
			return
		}
	}

	var user *database.User
	if query.WithOrders {
		user, err = h.db.GetUserWithOrders(id)
	} else {
		user, err = h.db.GetUser(id)
	}

	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot process request to db: " + err.Error()})
		}
		return
	}

	if data, err := json.Marshal(user); err == nil {
		if err := h.cacher.Set(ctx, cacheKey, string(data), userCacheTTL); err != nil {
			log.Printf("Cannot cache user with key: %q: %v", cacheKey, err)
		}
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUsers(c *gin.Context) {
	pageParam := c.DefaultQuery("page", "1")
	limitParam := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page: " + err.Error()})
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit: " + err.Error()})
		return
	}

	if page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "non positive page"})
		return
	}

	if limit <= 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 100"})
		return
	}

	resp, err := h.db.GetUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not execute query: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format: " + err.Error()})
		return
	}

	newUser, err := h.db.CreateUser(req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount: " + err.Error()})
		return
	}

	paramID := c.Param("id")
	id, err := utils.ValidateID(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.db.CreateOrder(uint(id), req.Amount)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot execute query: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) GetUserOrders(c *gin.Context) {
	paramID := c.Param("id")
	id, err := utils.ValidateID(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.db.GetUserOrders(id)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot execute query: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	paramID := c.Param("id")
	id, err := utils.ValidateID(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.db.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot execute query: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	cacheKey := fmt.Sprintf("user:%d", updatedUser.ID)
	if err := h.cacher.Del(ctx, cacheKey); err != nil {
		log.Printf("cannot delete key: %s from cache: %v", cacheKey, err)
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	paramID := c.Param("id")
	id, err := utils.ValidateID(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.DeleteUser(id); err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot execute query: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	cacheKey := fmt.Sprintf("user:%d", id)

	if err := h.cacher.Del(ctx, cacheKey); err != nil {
		log.Printf("cannot delete key: %s from cache: %v", cacheKey, err)
	}

	c.JSON(http.StatusNoContent, nil)
}
