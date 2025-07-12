package server

import (
	"context"
	"log/slog"

	v1 "github.com/ashtonx86/mocker/internal/handlers/v1"
	"github.com/ashtonx86/mocker/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

type WebServer struct {
	App *fiber.App
	Supervisor *supervisor.Supervisor
}

func NewWebServer() *WebServer {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() 

	app := fiber.New()
	
	su, err := supervisor.New(ctx)
	if err != nil {
		slog.Error("[func NewWebServer] init failed ", "error", err)
		return nil
	}

	su.Init()

	v1.MapV1Handlers(app, su)
	
	return &WebServer{
		App: app,
		Supervisor: su,
	}
}