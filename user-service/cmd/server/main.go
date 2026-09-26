package main

import (
	"log"
	user "microservices/proto/user"
	"net"
	"os"
	"userservice/internal/cache"
	"userservice/internal/database"
	grpcserver "userservice/internal/grpc"
	"userservice/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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

	// gRPC initializing

	grpcServer := grpc.NewServer()

	userServer := grpcserver.NewServer(db)

	user.RegisterUserServiceServer(
		grpcServer,
		userServer,
	)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		log.Println("gPRC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()

	// GIN initializing

	router := gin.Default()
	h := handlers.NewHandler(db, cacher)
	h.RegisterRouters(router)

	if err := router.Run(":8888"); err != nil {
		log.Fatalf("Cannot run server: %v", err)
	}
}
