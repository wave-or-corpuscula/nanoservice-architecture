package main

import (
	"log"
	"os"
	"userservice/internal/cache"
	"userservice/internal/database"
	"userservice/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatalln("DB_HOST variable is not set")
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatalln("Cannot initialize DB: ", err)
	}
	defer db.Close()

	cacher, err := cache.InitRedis()
	if err != nil {
		log.Println("Cannot initialize redis client:", err)
	} else {
		defer cacher.Close()
	}

	router := gin.Default()
	h := handlers.NewHandler(db, cacher)

	h.RegisterRouters(router)
	log.Println("Router started on :8888")
	if err := router.Run(":8888"); err != nil {
		log.Fatalf("Cannot run server: %v", err)
	}
}
