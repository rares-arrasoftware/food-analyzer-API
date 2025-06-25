package server

import "github.com/rares-arrasoftware/food-analyzer-api/v1/config"

// Server defines the interface for running the HTTP API server.
type Server interface {
	// Start launches the web server using the provided configuration.
	Start()
}

func NewServer(cfg config.Config) Server {
	// here we can decide if we want gin/fiber/custom etc.
	return newFiberServer(cfg)
}
