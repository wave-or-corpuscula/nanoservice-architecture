package cache

import (
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

func getTestClient() *RedisClient {
	cacher, err := InitRedis()
	if err != nil {
		log.Fatalln(err)
	}
	return cacher
}

func TestInitRedis(t *testing.T) {
	_ = getTestClient()
}

func TestGetSetDel(t *testing.T) {
	redis := getTestClient()

	key := "test:key"
	value := "test:value"

	err := redis.Set(t.Context(), key, value, time.Hour)
	assert.NoError(t, err)

	gotVal, err := redis.Get(t.Context(), key)
	assert.NoError(t, err)
	assert.Equal(t, value, gotVal)

	err = redis.Del(t.Context(), key)
	assert.NoError(t, err)

	_, err = redis.Get(t.Context(), key)
	assert.ErrorIs(t, err, ErrCacheMiss)
}
