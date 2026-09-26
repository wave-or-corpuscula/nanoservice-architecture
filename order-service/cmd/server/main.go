package main

import (
	"log"
	"orderservice/internal/database"
	"orderservice/internal/handlers"
	"os"

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

	router := gin.Default()
	h := handlers.NewHandler(db)
	h.RegisterRouters(router)
	log.Println("Router started on :9999")

	if err := router.Run(":9999"); err != nil {
		log.Fatalf("cannot run order-service: %v", err)
	}
}
