package mock

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ashtonx86/mocker/internal/data"
	"github.com/ashtonx86/mocker/internal/entities"
	"github.com/ashtonx86/mocker/internal/errs"
	"github.com/ashtonx86/mocker/internal/schemas"
	"github.com/google/uuid"
)
func CreateMock(ctx context.Context, db *sql.DB, mockData schemas.MockCreateRequest) (*entities.Mock, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, errs.NewError(err, errs.SQLErrorType, errs.ErrInternalFailure)
	}
	defer tx.Rollback()

	entity := entities.Mock{
		ID:           uuid.NewString(),
		Topic:        mockData.Topic,
		Instructions: mockData.Instructions,
		TimeMins:     mockData.TimeMins,
		AuthorID:     mockData.AuthorID,
		CreatedAt:    time.Now(),
		LastUpdatedAt: time.Now(),
	}

	cols := []string{"id", "topic", "instructions", "timeMins", "authorID", "createdAt", "lastUpdatedAt"}
	placeholders := make([]string, len(cols))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	stmt := fmt.Sprintf(`INSERT INTO mock (%s) VALUES (%s)`, strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	vals := []any{entity.ID, entity.Topic, entity.Instructions, entity.TimeMins, entity.AuthorID, entity.CreatedAt, entity.LastUpdatedAt}

	if _, err := tx.ExecContext(ctx, stmt, vals...); err != nil {
		return nil, data.SQLiteErrorComparator(err)
	}

	err = insertMockQuestions(ctx, tx, mockData, entity)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func insertMockQuestions(ctx context.Context, tx *sql.Tx, mockData schemas.MockCreateRequest, entity entities.Mock) (error) {
	mockQCols := []string{"id", "problem", "points", "correctOptionID", "mockID", "createdAt", "lastUpdatedAt"}
	mockQPlaceholders := make([]string, len(mockQCols))
	for i := range mockQPlaceholders {
		mockQPlaceholders[i] = "?"
	}

	mockQStmt := fmt.Sprintf(`INSERT INTO mockQuestion (%s) VALUES (%s)`, strings.Join(mockQCols, ", "), strings.Join(mockQPlaceholders, ", "))

	for _, q := range mockData.Questions {
		mockQ := entities.MockQuestion{
			ID:              uuid.NewString(),
			Problem:         q.Problem,
			Points:          q.Points,
			CorrectOptionID: q.CorrectOptionID,
			MockID:          entity.ID,
			CreatedAt:       time.Now(),
			LastUpdatedAt:   time.Now(),
		}

		mockQVals := []any{mockQ.ID, mockQ.Problem, mockQ.Points, mockQ.CorrectOptionID, mockQ.MockID, mockQ.CreatedAt, mockQ.LastUpdatedAt}
		if _, err := tx.ExecContext(ctx, mockQStmt, mockQVals...); err != nil {
			return data.SQLiteErrorComparator(err)
		}

		if err := insertMockOptions(ctx, q, mockQ.ID, tx); err != nil {
			return err
		}
	}
	return nil
}

func insertMockOptions(ctx context.Context, q schemas.MockQuestionSchema, questionID string, tx *sql.Tx) error {
	stmt := `INSERT INTO mockOption (id, number, option, questionID, createdAt, lastUpdatedAt) VALUES (?, ?, ?, ?, ?, ?)`

	for _, opt := range q.Options {
		option := entities.MockOption{
			ID:           uuid.NewString(),
			Number:       opt.Number,
			Option:       opt.Option,
			QuestionID:   questionID,
			CreatedAt:    time.Now(),
			LastUpdatedAt: time.Now(),
		}
		vals := []any{option.ID, option.Number, option.Option, option.QuestionID, option.CreatedAt, option.LastUpdatedAt}
		if _, err := tx.ExecContext(ctx, stmt, vals...); err != nil {
			return data.SQLiteErrorComparator(err)
		}
	}
	return nil
}
