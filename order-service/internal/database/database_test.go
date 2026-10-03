package database

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

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
			ExpectedError: ErrNotFound,
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
			ExpectedError: ErrNotFound,
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

func TestCreateEvent(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	order := Order{
		ID:        100,
		UserID:    10,
		Amount:    200,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	payload, err := json.Marshal(order)
	if err != nil {
		t.Fatal(err)
	}

	event, err := db.CreateEvent(t.Context(), "orders.created", "123", payload)
	assert.NoError(t, err)

	t.Log(event)
}

func TestGetUnpublished(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	eventType := "orders.created"

	orders, err := db.GetUnpublished(t.Context(), eventType)
	if err != nil {
		t.Fatal(err)
	}

	for _, order := range orders {
		t.Log(order)
	}
}

func TestSetEventPublished(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	err = db.SetEventPublished(t.Context(), 1)
	assert.NoError(t, err)
}
