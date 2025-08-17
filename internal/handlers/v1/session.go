package v1

import (
	"errors"

	"github.com/ashtonx86/mocker/internal/auth"
	"github.com/ashtonx86/mocker/internal/data"
	"github.com/ashtonx86/mocker/internal/errs"
	"github.com/ashtonx86/mocker/internal/schemas"
	"github.com/ashtonx86/mocker/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

// assert: SessionHandler implements Handler interface.
var _ Handler = (*SessionHandler)(nil)

type SessionHandler struct {
	Supervisor *supervisor.Supervisor
	SQLite     *data.SQLite
}
func NewSessionHandler(su *supervisor.Supervisor) *SessionHandler {
    return &SessionHandler{
        Supervisor: su,
        SQLite:     su.SQLite,
    }
}

func (h *SessionHandler) MapRoutes(router *fiber.Group) {
    router.Post("/", h.handlePOST)               
    router.Post("/answer", h.handleAddAnswer)    
    router.Get("/submit/:userID", h.handleSubmit)
}


func (h *SessionHandler) handlePOST(c *fiber.Ctx) error {
    user := auth.GetCurrentUser(c)
    if user == nil {
        return c.SendStatus(fiber.StatusInternalServerError)
    }

	req := new(schemas.SessionCreateRequest)
	c.BodyParser(&req)

	err := errs.Validate(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Bad request"))
	}

    sessions, err := h.Supervisor.SessionManager.New(c.Context(), req.MockID, user.ID)
    if err != nil {
        var e errs.Error
        if errors.As(err, &e) {
            switch e.Code {
            case errs.ErrNotFound:
                return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Mock with that ID does not exist"))
            case errs.ErrAlreadyExists:
                return c.Status(fiber.StatusForbidden).JSON(schemas.NewErrorAPIResponse(err, "Session for ths user already exists"))
            default:
                return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Internal error"))
            }
        }
    }

    return c.Status(fiber.StatusCreated).JSON(schemas.NewAPIResponse(true, sessions, ""))
}

func (h *SessionHandler) handleAddAnswer(c *fiber.Ctx) error {
    user := auth.GetCurrentUser(c)
    if user != nil {
        return c.SendStatus(fiber.StatusInternalServerError)
    }

	req := new(schemas.AnswerAddRequest)
	c.BodyParser(&req)

	err := errs.Validate(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Bad request"))
	}


    err = h.Supervisor.SessionManager.AddAnswer(c.Context(), req.MockID, user.ID, req.QuestionID, req.OptionID)

    if err != nil {
        var e errs.Error
        if errors.As(err, &e) {
            switch e.Code {
            case errs.ErrNotFound:
                return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Mock with that ID does not exist"))
            case errs.ErrAlreadyExists:
                return c.Status(fiber.StatusForbidden).JSON(schemas.NewErrorAPIResponse(err, "Session for ths user already exists"))
            default:
                return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Internal error"))
            }
        }
    }

    return c.SendStatus(fiber.StatusNoContent)
}

func (h *SessionHandler) handleSubmit(c *fiber.Ctx) error {
    userID := c.Params("userID")
    if userID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(errs.GenericBadRequstErr("user_id"), "Bad request"))
    }

    mockID := c.Query("mock_id")
    if mockID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(errs.GenericBadRequstErr("mock_id"), "Bad request"))
    }

    total, err := h.Supervisor.SessionManager.CalculateTotalMarks(c.Context(), h.SQLite.DB, mockID, userID)
    if err != nil {
        var e errs.Error
        if errors.As(err, &e) {
            switch e.Code {
            case errs.ErrNotFound:
                return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Mock with that ID does not exist"))
            case errs.ErrAlreadyExists:
                return c.Status(fiber.StatusForbidden).JSON(schemas.NewErrorAPIResponse(err, "Session for ths user already exists"))
            default:
                return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Internal error"))
            }
        }
    }

    return c.JSON(schemas.NewAPIResponse(true, fiber.Map{
        "user_id":   userID,
        "mock_id":   mockID,
        "total_marks": total,
    }, ""))
}
