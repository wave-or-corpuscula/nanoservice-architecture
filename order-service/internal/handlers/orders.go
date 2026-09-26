package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"orderservice/internal/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateID(paramID string) (uint, error) {
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	if id <= 0 {
		return 0, fmt.Errorf("non positive id")
	}

	return uint(id), nil
}

type CreateOrderRequest struct {
	UserID uint    `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data: " + err.Error()})
		return
	}

	// paramID := c.Param("id")
	// id, err := ValidateID(paramID)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	order, err := h.db.CreateOrder(req.UserID, req.Amount)
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
	id, err := ValidateID(paramID)
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
