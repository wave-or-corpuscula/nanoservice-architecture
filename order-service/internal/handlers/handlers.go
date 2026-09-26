package handlers

import (
	"net/http"
	"orderservice/internal/database"

	"github.com/gin-gonic/gin"
)

type Database interface {
	CreateOrder(userID uint, amount float64) (*database.Order, error)
	GetUserOrders(userID uint) (*database.OrdersResponse, error)
}

type Handler struct {
	db Database
}

func NewHandler(db Database) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) RegisterRouters(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/orders", h.CreateOrder)
		v1.GET("/users/:id/orders", h.GetUserOrders)
	}
}
