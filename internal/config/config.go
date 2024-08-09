package config

import "os"

type Config struct {
	ServerAddress string
	JWTSecret     string
	DBConnection  string
}

func NewConfig() *Config {
	serverAddress := os.Getenv("SERVER_ADDRESS")

	jwtSecret := os.Getenv("JWT_SECRET")

	dbConnection := os.Getenv("DB_CONNECTION")

	return &Config{
		ServerAddress: serverAddress,
		JWTSecret:     jwtSecret,
		DBConnection:  dbConnection,
	}
}
