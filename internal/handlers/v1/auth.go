package v1

import (
	"github.com/ashtonx86/mocker/internal/data"
	"github.com/ashtonx86/mocker/internal/schemas"
	"github.com/ashtonx86/mocker/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

// assert: AuthHandler implements Handler interface.
var _ Handler = (*AuthHandler)(nil)

type AuthHandler struct {
	Supervisor *supervisor.Supervisor

	SQLite *data.SQLite
}

func NewAuthHandler(su *supervisor.Supervisor) *AuthHandler {
	return &AuthHandler{
		Supervisor: su,

		SQLite: su.SQLite,
	}
}

func (h *AuthHandler) MapRoutes(router *fiber.Group) {
	router.Get("/", h.handleGET)
}

func (h *AuthHandler) handleGET(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(schemas.NewAPIResponse(true, "Hello world!", ""))
}