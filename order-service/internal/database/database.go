package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

var ErrUserNotFound = errors.New("user not found")

func GetEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func InitDB() (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	retryes := GetEnvAsInt("DB_CONNECTION_RETRYES", 5)
	timeout := time.Duration(GetEnvAsInt("DB_CONNECTION_RETRY_TIMEOUT", 2)) * time.Second

	var database *gorm.DB
	var err error

	for i := range retryes {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})

		if err == nil {
			break
		}

		log.Printf("Could not connect to DB, retrying %d/%d...\n", i+1, retryes)
		time.Sleep(timeout)
	}

	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("cannot configure connections pool: %w", err)
	}

	sqlDB.SetMaxIdleConns(GetEnvAsInt("DB_MAX_IDLE_CONNS", 10))
	sqlDB.SetMaxOpenConns(GetEnvAsInt("DB_MAX_OPEN_CONNS", 100))
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := database.AutoMigrate(&Order{}); err != nil {
		return nil, fmt.Errorf("cannot migrate Order table: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return &Database{
		db: database,
	}, nil
}

func (db *Database) Close() {
	sqlDB, err := db.db.DB()
	if err != nil {
		panic(err)
	}
	err = sqlDB.Close()
	if err != nil {
		panic(err)
	}
}

func (db *Database) CreateOrder(userID uint, amount float64) (*Order, error) {
	order := &Order{
		UserID: userID,
		Amount: amount,
	}

	if err := db.db.Debug().Create(order).Error; err != nil {
		// var pgErr *pgconn.PgError
		// if errors.As(err, &pgErr) && pgErr.Code == "23503" { // 23503 - code for foreign_key_violation in Postgres
		// 	return nil, ErrUserNotFound
		// }

		return nil, err
	}

	return order, nil
}

func (db *Database) GetUserOrders(userID uint) (*OrdersResponse, error) {
	var orders []Order
	if err := db.db.Debug().Where("user_id = ?", userID).First(&orders).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, ErrUserNotFound
		// }
		return nil, err
	}

	return &OrdersResponse{
		Orders: orders,
	}, nil
}
