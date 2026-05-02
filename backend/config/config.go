package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret     string
	JWTExpireHours int

	ServerPort string
}

var AppConfig *Config

func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	expireHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if expireHours == 0 {
		expireHours = 24
	}

	AppConfig = &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpireHours: expireHours,

		ServerPort: os.Getenv("SERVER_PORT"),
	}

	return nil
}
