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
