package main

import (
	"log"
	"orderservice/internal/database"
	"orderservice/internal/handlers"
	"os"

	user "microservices/proto/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// gRPC connection initializing

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	log.Println("gRPC server for order-service started")

	client := user.NewUserServiceClient(conn)

	// GIN initializing

	router := gin.Default()

	h := handlers.NewHandler(db, client)
	h.RegisterRouters(router)

	log.Println("Router started on :9999")

	if err := router.Run(":9999"); err != nil {
		log.Fatalf("cannot run order-service: %v", err)
	}
}
