package database

import (
	"os"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4/testutils/assert"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestCreateEvent(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	event := ProcessedEvent{
		ID:          1,
		ProcessedAt: time.Now(),
	}

	err = db.ProcessEvent(t.Context(), event.ID, event.ProcessedAt)
	assert.NoError(t, err)
}

func TestEventExists(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatal(err)
	}

	exists, err := db.EventExitsts(t.Context(), 1)
	assert.NoError(t, err)
	assert.Equal(t, exists, true)

	exists, err = db.EventExitsts(t.Context(), 100)
	assert.NoError(t, err)
	assert.Equal(t, exists, false)
}
