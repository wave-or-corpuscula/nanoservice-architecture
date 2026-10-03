package service

import (
	"log"
	"orderservice/internal/database"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalln("cannot load env variables: ", err)
	}

	code := m.Run()

	os.Exit(code)
}

func getOrdersService() *OrderService {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalln(err)
	}

	return NewOrderService(db)
}

func TestCreateOrder(t *testing.T) {
	var userID uint = 10
	var amount float64 = 399.99

	service := getOrdersService()
	order, err := service.CreateOrder(t.Context(), userID, amount)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(order)
}
