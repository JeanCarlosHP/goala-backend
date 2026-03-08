package domain

import "github.com/gofiber/fiber/v3"

type HTTPServer interface {
	StartServer() error
	ShutdownServer() error
	GetApp() *fiber.App
	GetLogger() Logger
	ConfigureMiddlewares()
}

type ListResponse struct {
	Items any `json:"items"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
