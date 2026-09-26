package main

import (
	"context"
	"fmt"
	"log"
	user "microservices/proto/user"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"userservice/internal/cache"
	"userservice/internal/config"
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

	cfg := config.Load()

	// gRPC initializing

	grpcServer := grpc.NewServer()

	userServer := grpcserver.NewServer(db)

	user.RegisterUserServiceServer(
		grpcServer,
		userServer,
	)

	grpcPort := os.Getenv("GRPC_PORT")
	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%s", grpcPort),
	)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		log.Printf("gPRC server listening on :%s\n", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()

	// GIN initializing

	router := gin.Default()
	h := handlers.NewHandler(db, cacher, cfg)
	h.RegisterRouters(router)

	srv := http.Server{
		Addr:    ":8888",
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Cannot run server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTTIN, syscall.SIGTERM)

	<-quit
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("GIN server forced to shutdown: %v\n", err)
	} else {
		log.Println("GIN server stopped gracefully")
	}
	log.Println("Application exiting successfully")
}
