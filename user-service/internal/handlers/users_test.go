package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"userservice/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockDatabase struct{}

func getTestUser(name string, email string) *database.User {
	return &database.User{
		ID:        1,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (md *MockDatabase) CreateUser(name string, email string) (*database.User, error) {
	return getTestUser(name, email), nil
}

func (md *MockDatabase) GetUser(id uint) (*database.User, error) {
	user := getTestUser("test", "test")
	user.ID = id
	return user, nil
}

func (md *MockDatabase) GetUsers(page int, limit int) (*database.UsersPaginationResponse, error) {
	return &database.UsersPaginationResponse{
		Users: []database.User{
			*getTestUser("test1", "email1"),
			*getTestUser("test2", "email2"),
			*getTestUser("test3", "email3"),
		},
		Total: 3,
		Page:  page,
		Limit: limit,
	}, nil
}

func (md *MockDatabase) GetUserWithOrders(id uint) (*database.User, error) {
	user := getTestUser("test", "test")
	user.ID = id
	return user, nil
}

func (md *MockDatabase) UpdateUser(userID uint, name string, email string) (*database.User, error) {
	user := getTestUser(name, email)
	user.ID = userID
	return user, nil
}

func (md *MockDatabase) DeleteUser(userID uint) error {
	return nil
}

type MockCacher struct{}

var ErrMockCacher = errors.New("mock cache")

func (mc *MockCacher) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return ErrMockCacher
}

func (mc *MockCacher) Get(ctx context.Context, key string) (string, error) {
	return "", ErrMockCacher
}

func (mc *MockCacher) Del(ctx context.Context, keys ...string) error {
	return ErrMockCacher
}

func (mc *MockCacher) Close() error {
	return ErrMockCacher
}

func getTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	db := &MockDatabase{}
	cacher := &MockCacher{}

	router := gin.New()
	handlers := NewHandler(db, cacher)
	handlers.RegisterRouters(router)

	return router
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name             string
		UserID           string
		ExpectedCode     int
		ExpectedErrorStr string
	}{
		{
			name:             "success",
			UserID:           "1",
			ExpectedCode:     http.StatusOK,
			ExpectedErrorStr: "",
		},
	}

	r := getTestRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/v1/users/%s", tt.UserID),
				nil,
			)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.ExpectedCode, rec.Code)
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name             string
		userName         string
		userEmail        string
		ExpectedCode     int
		ExpectedErrorStr string
	}{
		{
			name:             "success",
			userName:         "name",
			userEmail:        "email",
			ExpectedCode:     http.StatusCreated,
			ExpectedErrorStr: "",
		},
		{
			name:             "fail - empty name",
			userName:         "",
			userEmail:        "email",
			ExpectedCode:     http.StatusBadRequest,
			ExpectedErrorStr: "Field validation for 'Name'",
		},
		{
			name:             "fail - empty email",
			userName:         "name",
			userEmail:        "",
			ExpectedCode:     http.StatusBadRequest,
			ExpectedErrorStr: "Field validation for 'Email'",
		},
	}

	r := getTestRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := fmt.Sprintf(
				"{\"name\":\"%s\", \"email\":\"%s\"}",
				tt.userName,
				tt.userEmail,
			)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/users",
				strings.NewReader(body),
			)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.ExpectedCode, rec.Code)

			if tt.ExpectedErrorStr != "" {
				errorStruct := struct {
					Error string `json:"error"`
				}{}

				err := json.NewDecoder(rec.Body).Decode(&errorStruct)
				if err != nil {
					t.Fatal(err)
				}

				assert.Contains(t, errorStruct.Error, tt.ExpectedErrorStr)
			}
		})
	}
}
