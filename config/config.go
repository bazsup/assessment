package config

import "os"

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable: " + name)
	}
	return v
}

// Config contains app config like running port and database url
type Config struct {
	Port        string
	DatabaseUrl string
}

// NewConfig returns a new app config
func NewConfig() *Config {
	return &Config{Port: getenv("PORT"), DatabaseUrl: getenv("DATABASE_URL")}
}
