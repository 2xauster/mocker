package v1

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/ashtonx86/mocker/internal/auth"
	"github.com/ashtonx86/mocker/internal/data"
	"github.com/ashtonx86/mocker/internal/errs"
	"github.com/ashtonx86/mocker/internal/logging"
	"github.com/ashtonx86/mocker/internal/mock"
	"github.com/ashtonx86/mocker/internal/schemas"
	"github.com/ashtonx86/mocker/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

const MOCK_TIMEOUT = 40 * time.Second

// assert: MockHandler implements Handler interface.
var _ Handler = (*MockHandler)(nil)

type MockHandler struct {
	Supervisor *supervisor.Supervisor
	SQLite     *data.SQLite
}

func NewMockHandler(su *supervisor.Supervisor) *MockHandler {
	return &MockHandler{
		Supervisor: su,
		SQLite:     su.SQLite,
	}
}

func (h *MockHandler) MapRoutes(router *fiber.Group) {
	router.Post("/", h.handlePOST) 
	router.Get("/:id", h.handleGET)
}

func (h *MockHandler) handlePOST(c *fiber.Ctx) error {
	user := auth.GetCurrentUser(c)

	if user == nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	req := new(schemas.MockCreateRequest)
	c.BodyParser(&req)

	err := errs.Validate(req)
	if err != nil && errors.Is(err, errs.Error{Code: errs.ErrDataIllegal, Type: errs.DataErrorType.String()}) {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Bad request"))
	}

	req.AuthorID = user.ID 
	
	ctx, cancel := context.WithTimeout(context.Background(), MOCK_TIMEOUT)
	defer cancel()

	entity, err := mock.CreateMock(ctx, h.SQLite.DB, *req)
	if err != nil {
		var e errs.Error
		if errors.As(err, &e) {
			logging.Log(slog.LevelError, c, "Mock creation failed", "user_id", req.AuthorID, "error", e)

			switch e.Code {
			case errs.ErrAlreadyExists:
				return c.Status(fiber.StatusForbidden).JSON(schemas.NewErrorAPIResponse(err, "Already exists"))
			case errs.ErrDataMismatch:
				return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Data mismatch"))
			case errs.ErrInternalFailure:
				return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Internal failure"))
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Unknown error"))
	}

	return c.JSON(schemas.NewAPIResponse(true, entity, ""))
}

func (h *MockHandler) handleGET(c *fiber.Ctx) error {
	mockID := c.Params("id")
	if mockID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(errors.New("missing mock ID"), "Bad request"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), MOCK_TIMEOUT)
	defer cancel()

	entity, err := mock.GetMock(ctx, h.SQLite.DB, mockID)
	if err != nil {
		var e errs.Error
		if errors.As(err, &e) {
			slog.Error("[pkg handlers : mock.go : func handleGET] Failed to fetch mock", "error", err)

			switch e.Code {
			case errs.ErrNotFound:
				return c.Status(fiber.StatusNotFound).JSON(schemas.NewErrorAPIResponse(err, "Not found"))
			case errs.ErrInternalFailure:
				return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Internal failure"))
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Unknown error"))
	}

	return c.JSON(schemas.NewAPIResponse(true, entity, ""))
}
