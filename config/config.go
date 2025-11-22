package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	PostgresDSN     string
	MongoURI        string
	MongoDB         string
	JWTSecret       string
	JWTExpiresHours int
}

var C Config

func Load() {
	_ = godotenv.Load()

	C = Config{
		AppPort:         get("APP_PORT", "3000"),
		PostgresDSN:     get("POSTGRES_DSN", ""),
		MongoURI:        get("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:         get("MONGO_DB", "prestasi_mhs"),
		JWTSecret:       get("JWT_SECRET", "secret123"),
		JWTExpiresHours: getInt("JWT_EXPIRES_HOURS", 24),
	}
}

func get(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
