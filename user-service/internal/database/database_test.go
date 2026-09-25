package database

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalln("cannot load env variables: ", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestGetUsers(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	resp, _ := db.GetUsers(3, 5)
	t.Log("Page:", resp.Page)
	t.Log("Limit:", resp.Limit)
	t.Log("Total:", resp.Total)
	for _, user := range resp.Users {
		t.Log(user)
	}
}

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		Name          string
		UserID        uint
		Amount        float64
		ExpectedError error
	}{
		{
			Name:          "success",
			UserID:        1,
			Amount:        1000.1,
			ExpectedError: nil,
		},
		{
			Name:          "unknown user",
			UserID:        0,
			Amount:        1000,
			ExpectedError: ErrUserNotFound,
		},
		{
			Name:          "negative amount",
			UserID:        1,
			Amount:        -1000,
			ExpectedError: nil,
		},
	}

	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			order, err := db.CreateOrder(tt.UserID, tt.Amount)
			assert.ErrorIs(t, err, tt.ExpectedError)
			t.Log(order)
		})
	}
}

func TestGetUserOrders(t *testing.T) {
	tests := []struct {
		Name          string
		UserID        uint
		ExpectedError error
	}{
		{
			Name:          "success",
			UserID:        1,
			ExpectedError: nil,
		},
		{
			Name:          "unknown user",
			UserID:        0,
			ExpectedError: ErrUserNotFound,
		},
	}

	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			resp, err := db.GetUserOrders(tt.UserID)
			assert.ErrorIs(t, err, tt.ExpectedError)

			t.Log(resp)
		})
	}
}

func TestGetUserWithOrders(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	user, err := db.GetUserWithOrders(1)
	assert.NoError(t, err)

	t.Log(user)
}

func TestUpdateUser(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	user, err := db.GetUser(1)
	assert.NoError(t, err)

	user.Name += "1"
	user.Email += "1"

	updatedUser, err := db.UpdateUser(user.ID, user.Name, user.Email)
	assert.NoError(t, err)

	assert.Equal(t, user.Name, updatedUser.Name)
	assert.Equal(t, user.Email, updatedUser.Email)
}

func TestUpdateUnknownUser(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	user, err := db.GetUser(1)
	assert.NoError(t, err)

	user.Name += "1"
	user.Email += "1"

	updatedUser, err := db.UpdateUser(user.ID, user.Name, user.Email)
	assert.NoError(t, err)

	assert.Equal(t, user.Name, updatedUser.Name)
	assert.Equal(t, user.Email, updatedUser.Email)
}

func TestDeleteUser(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	err = db.DeleteUser(999)
	assert.ErrorIs(t, err, ErrUserNotFound)

	newUser, err := db.CreateUser("test", "test")
	assert.NoError(t, err)

	err = db.DeleteUser(newUser.ID)
	assert.NoError(t, err)

	_, err = db.GetUser(newUser.ID)
	assert.ErrorIs(t, err, ErrUserNotFound)
}
