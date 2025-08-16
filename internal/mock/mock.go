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
	"github.com/ashtonx86/mocker/internal/utils"
	"github.com/google/uuid"
)

func CreateMock(ctx context.Context, db *sql.DB, mockData schemas.MockCreateRequest) (*entities.Mock, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, errs.NewError(err, errs.SQLErrorType, errs.ErrInternalFailure)
	}
	defer tx.Rollback()

	entity := entities.Mock{
		ID:            uuid.NewString(),
		Topic:         mockData.Topic,
		Instructions:  mockData.Instructions,
		TimeMins:      mockData.TimeMins,
		AuthorID:      mockData.AuthorID,
		CreatedAt:     time.Now(),
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

type FullMock struct {
	entities.Mock
	Questions []FullMockQuestion `json:"questions"`
}

type FullMockQuestion struct {
	entities.MockQuestion
	Options []entities.MockOption `json:"options"`
}

func GetMock(ctx context.Context, db *sql.DB, id string) (*FullMock, error) {
	var mock entities.Mock
	mockStmt := `
        SELECT id, topic, instructions, timeMins, authorID, createdAt, lastUpdatedAt
        FROM mock
        WHERE id = ?
    `

	var createdAtString, lastUpdatedAtString string 

	err := db.QueryRowContext(ctx, mockStmt, id).Scan(
		&mock.ID,
		&mock.Topic,
		&mock.Instructions,
		&mock.TimeMins,
		&mock.AuthorID,
		&createdAtString,
		&lastUpdatedAtString,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewError(err, errs.DataErrorType, errs.ErrNotFound)
		}
		return nil, data.SQLiteErrorComparator(err)
	}

	createdAt, _ := utils.ParseTime(createdAtString)
	lastUpdatedAt, _ := utils.ParseTime(lastUpdatedAtString)

	mock.CreatedAt = *createdAt
	mock.LastUpdatedAt = *lastUpdatedAt

	qStmt := `
        SELECT id, problem, points, mockID, createdAt, lastUpdatedAt
        FROM mockQuestion
        WHERE mockID = ?
    `
	rows, err := db.QueryContext(ctx, qStmt, mock.ID)
	if err != nil {
		return nil, data.SQLiteErrorComparator(err)
	}
	defer rows.Close()

	var fullQuestions []FullMockQuestion
	for rows.Next() {
		var q entities.MockQuestion

		var qCreatedAtStr, qLastUpdatedAtStr string 

		if err := rows.Scan(
			&q.ID,
			&q.Problem,
			&q.Points,
			&q.MockID,
			&qCreatedAtStr,
			&qLastUpdatedAtStr,
		); err != nil {
			return nil, err
		}

		qCreatedAt, _ := utils.ParseTime(qCreatedAtStr)
		qLastUpdatedAt, _ := utils.ParseTime(qLastUpdatedAtStr)

		q.CreatedAt = *qCreatedAt
		q.LastUpdatedAt = *qLastUpdatedAt

		optStmt := `
            SELECT id, number, option, questionID, createdAt, lastUpdatedAt
            FROM mockOption
            WHERE questionID = ?
        `
		optRows, err := db.QueryContext(ctx, optStmt, q.ID)
		if err != nil {
			return nil, data.SQLiteErrorComparator(err)
		}

		var options []entities.MockOption
		for optRows.Next() {
			var opt entities.MockOption

			var optCreatedAtStr, optLastUpdatedAtStr string 

			if err := optRows.Scan(
				&opt.ID,
				&opt.Number,
				&opt.Option,
				&opt.QuestionID,
				&optCreatedAtStr,
				&optLastUpdatedAtStr,
			); err != nil {
				optRows.Close()
				return nil, err
			}

			optCreatedAt, _ := utils.ParseTime(optCreatedAtStr)
			optLastUpdatedAt, _ := utils.ParseTime(optLastUpdatedAtStr)

			opt.CreatedAt = *optCreatedAt
			opt.LastUpdatedAt = *optLastUpdatedAt

			options = append(options, opt)
		}
		optRows.Close()

		fullQuestions = append(fullQuestions, FullMockQuestion{
			MockQuestion: q,
			Options:      options,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &FullMock{
		Mock:      mock,
		Questions: fullQuestions,
	}, nil
}

func insertMockQuestions(ctx context.Context, tx *sql.Tx, mockData schemas.MockCreateRequest, entity entities.Mock) error {
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
			ID:            uuid.NewString(),
			Number:        opt.Number,
			Option:        opt.Option,
			QuestionID:    questionID,
			CreatedAt:     time.Now(),
			LastUpdatedAt: time.Now(),
		}
		vals := []any{option.ID, option.Number, option.Option, option.QuestionID, option.CreatedAt, option.LastUpdatedAt}
		if _, err := tx.ExecContext(ctx, stmt, vals...); err != nil {
			return data.SQLiteErrorComparator(err)
		}
	}
	return nil
}
