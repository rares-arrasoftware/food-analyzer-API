package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rares-arrasoftware/food-analyzer-api/v1/config"
	"github.com/rares-arrasoftware/food-analyzer-api/v1/server"

	"github.com/spf13/cobra"
)

var (
	cfg = config.DefaultConfig()

	rootCmd = &cobra.Command{
		Use:   "food-analyzer",
		Short: "A simple API to analyze food images and manage users.",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Starting Food Analyzer API...")
			srv := server.NewServer(*cfg)
			srv.Start()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.Port, "port", cfg.Port, "Port to listen on for Fiber server")
	rootCmd.PersistentFlags().StringVar(&cfg.JWTSecret, "jwt-secret", cfg.JWTSecret, "JWT signing key")
	rootCmd.PersistentFlags().IntVar(&cfg.JWTExpiry, "jwt-expiry", cfg.JWTExpiry, "JWT token expiry time in hours")
	rootCmd.PersistentFlags().StringVar(&cfg.DatabaseDSN, "db", cfg.DatabaseDSN, "Database DSN")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
