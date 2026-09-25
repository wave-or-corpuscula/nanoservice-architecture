package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"userservice/pkg/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

var ErrUserNotFound = errors.New("user not found")

func InitDB() (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	retryes := utils.GetEnvAsInt("DB_CONNECTION_RETRYES", 5)
	timeout := time.Duration(utils.GetEnvAsInt("DB_CONNECTION_RETRY_TIMEOUT", 2)) * time.Second

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

	sqlDB.SetMaxIdleConns(utils.GetEnvAsInt("DB_MAX_IDLE_CONNS", 10))
	sqlDB.SetMaxOpenConns(utils.GetEnvAsInt("DB_MAX_OPEN_CONNS", 100))
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := database.AutoMigrate(&User{}); err != nil {
		return nil, fmt.Errorf("cannot migrate User table: %w", err)
	}

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

func (db *Database) CreateUser(name string, email string) (*User, error) {
	user := &User{Name: name, Email: email}
	if err := db.db.Debug().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (db *Database) GetUser(id uint) (*User, error) {
	var user User

	if err := db.db.Debug().First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db *Database) GetUserWithOrders(id uint) (*User, error) {
	var user User

	if err := db.db.Debug().Preload("Orders").First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) GetUsers(page int, limit int) (*UsersPaginationResponse, error) {
	var users []User
	var total int64

	if err := db.db.Debug().Model(&User{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := limit * (page - 1)
	if err := db.db.Debug().Limit(limit).Offset(offset).Order("id ASC").Find(&users).Error; err != nil {
		return nil, err
	}

	resp := &UsersPaginationResponse{
		Users: users,
		Page:  page,
		Limit: limit,
		Total: int(total),
	}
	return resp, nil
}

func (db *Database) CreateOrder(userID uint, amount float64) (*Order, error) {
	order := &Order{
		UserID: userID,
		Amount: amount,
	}

	if err := db.db.Debug().Create(order).Error; err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // 23503 - code for foreign_key_violation in Postgres
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return order, nil
}

func (db *Database) GetUserOrders(userID uint) (*OrdersResponse, error) {
	var user User

	if err := db.db.Debug().Preload("Orders").First(&user, userID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &OrdersResponse{
		Orders: user.Orders,
	}, nil
}

func (db *Database) UpdateUser(userID uint, name string, email string) (*User, error) {
	var updated User

	result := db.db.Debug().
		Model(&updated).
		Clauses(clause.Returning{}).
		Where("id = ?", userID).
		Updates(
			User{
				Name:  name,
				Email: email,
			},
		)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrUserNotFound
	}

	return &updated, nil
}

func (db *Database) DeleteUser(userID uint) error {
	result := db.db.Delete(&User{}, userID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
