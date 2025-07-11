package server

import "github.com/gofiber/fiber/v2"

type WebServer struct {
	App *fiber.App
}

func NewWebServer() *WebServer {
	return &WebServer{}
}