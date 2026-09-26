package handlers

import (
	"context"
	"time"
	"userservice/internal/database"

	"github.com/gin-gonic/gin"
)

type Database interface {
	CreateUser(name string, email string) (*database.User, error)
	GetUser(id uint) (*database.User, error)
	GetUsers(page int, limit int) (*database.UsersPaginationResponse, error)
	GetUserWithOrders(id uint) (*database.User, error)
	UpdateUser(userID uint, name string, email string) (*database.User, error)
	DeleteUser(userID uint) error
}

type Cacher interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Close() error
}

const userCacheTTL = 10 * time.Second

type Handler struct {
	db     Database
	cacher Cacher
}

func NewHandler(db Database, cacher Cacher) *Handler {
	return &Handler{
		db:     db,
		cacher: cacher,
	}
}

func (h *Handler) RegisterRouters(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/users/:id", h.GetUser)
		v1.PUT("/users/:id", h.UpdateUser)
		v1.DELETE("/users/:id", h.DeleteUser)
		v1.GET("/users", h.GetUsers)
		v1.POST("/users", h.CreateUser)
	}
}
