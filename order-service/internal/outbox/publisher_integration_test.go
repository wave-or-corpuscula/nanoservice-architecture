package outbox

import (
	"context"
	"log"
	"orderservice/internal/database"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type MockEventPublisher struct{}

func (mp *MockEventPublisher) Publish(ctx context.Context, key string, value []byte) error {
	return nil
}

func (mp *MockEventPublisher) PublishOrder(ctx context.Context, order database.Order) error {
	return nil
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalln("cannot load env variables: ", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestRun(t *testing.T) {
	db, err := database.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	p := &MockEventPublisher{}

	eventType := "orders.created"

	publisher := New(eventType, db, p)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	publisher.Run(ctx)
}
