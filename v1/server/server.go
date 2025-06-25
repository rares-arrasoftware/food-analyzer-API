package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rares-arrasoftware/food-analyzer-api/v1/analyzer"
	"github.com/rares-arrasoftware/food-analyzer-api/v1/auth"
	"github.com/rares-arrasoftware/food-analyzer-api/v1/config"
)

type fiberServer struct {
	app             *fiber.App
	config          config.Config
	authService     auth.Service
	analyzerService analyzer.Service
}

func newFiberServer(cfg config.Config) Server {
	app := fiber.New()

	s := &fiberServer{
		app:             app,
		config:          cfg,
		authService:     auth.NewService(cfg),
		analyzerService: analyzer.NewService(),
	}

	app.Static("/api", "docs/")

	// Register routes using method receivers
	app.Post("/auth/register", s.registerHandler())
	app.Post("/auth/login", s.loginHandler())

	// Protected routes
	protected := app.Group("/food", s.jwtMiddleware())
	protected.Post("/analyze", s.analyzeHandler())

	return s
}

func (s *fiberServer) Start() {
	log.Printf("Starting server on port %s...", s.config.Port)
	log.Fatal(s.app.Listen(":" + s.config.Port))
}

func (s *fiberServer) registerHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req auth.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}

		resp, err := s.authService.Register(req)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(resp)
	}
}

func (s *fiberServer) loginHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req auth.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}

		resp, err := s.authService.Login(req)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(resp)
	}
}

func (s *fiberServer) analyzeHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		resp := s.analyzerService.Analyze()
		return c.JSON(resp)
	}
}

func (s *fiberServer) jwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		claims, err := s.authService.ValidateTokenFromHeader(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		c.Locals("userID", claims.Sub)
		c.Locals("email", claims.Email)
		return c.Next()
	}
}
