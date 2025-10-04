package config

import "github.com/joho/godotenv"

func LoadConfig() error {
	return godotenv.Load()
}
