package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"orderservice/internal/database"
	"strconv"
	"strings"
	"testing"
	"time"

	user "microservices/proto/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func getTestOrder(userID uint, amount float64) *database.Order {
	return &database.Order{
		ID:        1,
		Amount:    amount,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type MockDatabase struct{}

func (md *MockDatabase) CreateOrder(userID uint, amount float64) (*database.Order, error) {
	return getTestOrder(userID, amount), nil
}

func (md *MockDatabase) GetUserOrders(userID uint) (*database.OrdersResponse, error) {
	return &database.OrdersResponse{
		Orders: []database.Order{
			*getTestOrder(userID, 100),
			*getTestOrder(userID, 200),
			*getTestOrder(userID, 300),
		},
	}, nil
}

type MockUserClient struct{}

func (mc *MockUserClient) GetUser(
	ctx context.Context,
	in *user.GetUserRequest,
	opts ...grpc.CallOption,
) (*user.GetUserResponse, error) {
	return &user.GetUserResponse{Id: 1, Name: "grpcName", Email: "grpcEmail"}, nil
}

func getTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	client := &MockUserClient{}
	db := &MockDatabase{}
	handler := NewHandler(db, client)

	router := gin.New()
	handler.RegisterRouters(router)

	return router
}

func TestPing(t *testing.T) {
	router := getTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	content, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "pong", string(content))
}

func TestCreateOrder(t *testing.T) {
	r := getTestRouter()

	tests := []struct {
		name             string
		UserID           string
		Amount           string
		ExpectedErrorStr string
		ExpectedCode     int
	}{
		{
			name:             "success",
			UserID:           "1",
			Amount:           "100",
			ExpectedErrorStr: "",
			ExpectedCode:     http.StatusCreated,
		},
		{
			name:             "success - float amount",
			UserID:           "1",
			Amount:           "100.38",
			ExpectedErrorStr: "",
			ExpectedCode:     http.StatusCreated,
		},
		{
			name:             "fail - 0 user id",
			UserID:           "0",
			Amount:           "100.38",
			ExpectedErrorStr: "failed on the 'required' tag",
			ExpectedCode:     http.StatusBadRequest,
		},
		{
			name:             "fail - negative user id",
			UserID:           "-1",
			Amount:           "100.38",
			ExpectedErrorStr: "cannot unmarshal number",
			ExpectedCode:     http.StatusBadRequest,
		},
		{
			name:             "fail - zero amount",
			UserID:           "1",
			Amount:           "0",
			ExpectedErrorStr: "failed on the 'required' tag",
			ExpectedCode:     http.StatusBadRequest,
		},
		{
			name:             "fail - negative amount",
			UserID:           "1",
			Amount:           "-10",
			ExpectedErrorStr: "failed on the 'gt' tag",
			ExpectedCode:     http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := fmt.Sprintf(
				"{\"user_id\":%s,\"amount\":%s}",
				tt.UserID,
				tt.Amount,
			)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/orders",
				strings.NewReader(body),
			)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.ExpectedCode, rec.Code)

			if tt.ExpectedErrorStr == "" {
				var newOrder database.Order
				err := json.NewDecoder(rec.Body).Decode(&newOrder)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.UserID, strconv.FormatUint(uint64(newOrder.UserID), 10))
				assert.Equal(t, tt.Amount, strconv.FormatFloat(newOrder.Amount, 'f', -1, 64))
			} else {
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

func TestGetUserOrders(t *testing.T) {
	tests := []struct {
		name             string
		UserID           string
		ExpectedErrorStr string
		ExpectedCode     int
	}{
		{
			name:             "success",
			UserID:           "1",
			ExpectedCode:     http.StatusOK,
			ExpectedErrorStr: "",
		},
		{
			name:             "fail - zero id",
			UserID:           "0",
			ExpectedCode:     http.StatusBadRequest,
			ExpectedErrorStr: "non positive id",
		},
		{
			name:             "fail - negative id",
			UserID:           "-10",
			ExpectedCode:     http.StatusBadRequest,
			ExpectedErrorStr: "non positive id",
		},
		{
			name:             "fail - invalid id",
			UserID:           "abc",
			ExpectedCode:     http.StatusBadRequest,
			ExpectedErrorStr: "invalid syntax",
		},
	}

	r := getTestRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/v1/users/%s/orders", tt.UserID),
				nil,
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
