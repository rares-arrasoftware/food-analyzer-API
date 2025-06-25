package config

// Package config defines the central application configuration.

// We use a plain struct instead of an interface
// to keep things simple. Interfaces can be added later if
// dynamic or mockable config becomes necessary.
type Config struct {
	Port        string
	JWTSecret   string
	JWTExpiry   int // in hours
	DatabaseDSN string
}

func DefaultConfig() *Config {
	return &Config{
		Port:        "3000",
		JWTSecret:   "defaultsecret",
		JWTExpiry:   72,
		DatabaseDSN: "food.db",
	}
}
