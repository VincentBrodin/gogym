package config

import (
	"github.com/joho/godotenv"
)

// JwtSecret holds the JWT secret key
var ENV map[string]string

func init() {
	env, _ := godotenv.Read(".env")

	ENV = env
}
