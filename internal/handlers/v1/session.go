package v1

import (
	"github.com/ashtonx86/mocker/internal/data"
	"github.com/ashtonx86/mocker/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

// assert: SessionHandler implements Handler interface.
var _ Handler = (*SessionHandler)(nil)

type SessionHandler struct {
	Supervisor *supervisor.Supervisor
	SQLite     *data.SQLite
}

func NewSessionHandler(su *supervisor.Supervisor) *MockHandler {
	return &MockHandler{
		Supervisor: su,
		SQLite:     su.SQLite,
	}
}

func (h *SessionHandler) MapRoutes(router *fiber.Group) {

}