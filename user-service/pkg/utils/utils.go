package utils

import (
	"fmt"
	"os"
	"strconv"
)

func ValidateID(paramID string) (uint, error) {
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err.Error())
	}

	if id <= 0 {
		return 0, fmt.Errorf("non positive id")
	}

	return uint(id), nil
}

func GetEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}
