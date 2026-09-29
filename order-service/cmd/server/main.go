package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"orderservice/internal/config"
	"orderservice/internal/database"
	"orderservice/internal/handlers"
	"orderservice/internal/kafka"
	"os"
	"os/signal"
	"syscall"

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

	cfg := config.Load()

	// Kafka initialization

	producer := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaClientID)
	defer producer.Close()

	// gRPC connection initializing

	grpcURL := fmt.Sprintf(
		"%s:%s",
		os.Getenv("GRPC_USER_HOST"),
		os.Getenv("GRPC_USER_PORT"),
	)

	conn, err := grpc.NewClient(
		grpcURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	log.Printf("gRPC connection started on %s\n", grpcURL)

	client := user.NewUserServiceClient(conn)

	// GIN initializing

	router := gin.Default()

	h := handlers.NewHandler(db, client, cfg, producer)
	h.RegisterRouters(router)

	srv := http.Server{
		Addr:    ":9999",
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("cannot run order-service: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("GIN server forced to shutdown: %v\n", err)
	} else {
		log.Println("GIN server stopped gracefully")
	}
	log.Println("Application exiting successfully")
}
