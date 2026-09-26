package handlers

import (
	"time"
	"userservice/internal/database"
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
