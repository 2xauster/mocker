package session

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ashtonx86/mocker/internal/data"
	"github.com/ashtonx86/mocker/internal/errs"
	"github.com/ashtonx86/mocker/internal/mock"
	"github.com/google/uuid"
)

type SessionManager struct {
	DB    *sql.DB
	Redis *data.Redis
}

type AnswerResult struct {
    QuestionID    string
    SelectedOption string
    IsCorrect     bool
}

func NewSessionManager(db *sql.DB, redisClient *data.Redis) *SessionManager {
	return &SessionManager{
		DB:    db,
		Redis: redisClient,
	}
}

// Create new session
func (s *SessionManager) New(ctx context.Context, mockID string, userID string) (map[string]Session, error) {
	ses := Session{
		ID:     uuid.NewString(),
		MockID: mockID,
		UserID: userID,

		CreatedAt: time.Now(),
	}
	identifier := userID

	d, err := mock.GetMock(ctx, s.DB, mockID)
	if err != nil {
		return nil, err
	}

	sesH, err := json.Marshal(&ses)
	if err != nil {
		return nil, errs.NewError(err, errs.DataErrorType, errs.ErrUndefined)
	}
	
	err = s.Redis.Client.Set(ctx, identifier, sesH, time.Duration(d.TimeMins)).Err()
	err = data.RedisErrorComparator(err) 

	if err != nil {
		return nil, err
	}

	return map[string]Session{
		identifier: ses,
	}, nil
}

func (s *SessionManager) AddAnswer(ctx context.Context, mockID string, userID string, questionID string, optionID string) error {
	b, err := s.Redis.Client.Get(ctx, userID).Bytes()
	err = data.RedisErrorComparator(err)

	if err != nil {
		return err
	}

	var ses Session
	if err = json.Unmarshal(b, &ses); err != nil {
		return errs.NewError(err, errs.DataErrorType, errs.ErrUndefined)
	}

	if ses.Answers == nil {
		ses.Answers = make(map[string]string)
	}

	ses.Answers[questionID] = optionID
	sesH, err := json.Marshal(&ses)
	if err != nil {
		return errs.NewError(err, errs.DataErrorType, errs.ErrUndefined)
	}

	err = s.Redis.Client.Set(ctx, userID, sesH, time.Duration(ses.TTL)*time.Second).Err()
	err = data.RedisErrorComparator(err)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionManager) CalculateTotalMarks(ctx context.Context, db *sql.DB, mockID string, userID string) (int, error) {
    mck, err := mock.GetMock(ctx, db, mockID)
    if err != nil {
        return 0, err
    }

    b, err := s.Redis.Client.Get(ctx, userID).Bytes()
    err = data.RedisErrorComparator(err)
    if err != nil {
        return 0, err
    }

    var ses Session
    if err = json.Unmarshal(b, &ses); err != nil {
        return 0, errs.NewError(err, errs.DataErrorType, errs.ErrUndefined)
    }

    if ses.Answers == nil {
        return 0, nil
    }

    total := 0
    for _, q := range mck.Questions {
        if optionID, ok := ses.Answers[q.ID]; ok {
            if optionID == q.CorrectOptionID {
                total += q.Points
            } else {
                total -= q.Points
            }
        }
    }
    return total, nil
}


func (s *SessionManager) GetAnswerResults(ctx context.Context, db *sql.DB, mockID string, userID string) ([]AnswerResult, error) {
    mck, err := mock.GetMock(ctx, db, mockID)
    if err != nil {
        return nil, err
    }

    b, err := s.Redis.Client.Get(ctx, userID).Bytes()
    err = data.RedisErrorComparator(err)
    if err != nil {
        return nil, err
    }

    var ses Session
    if err = json.Unmarshal(b, &ses); err != nil {
        return nil, errs.NewError(err, errs.DataErrorType, errs.ErrUndefined)
    }

    results := []AnswerResult{}

    if ses.Answers == nil {
        return results, nil 
    }

    for _, q := range mck.Questions {
        selectedOption, answered := ses.Answers[q.ID]
        isCorrect := answered && selectedOption == q.CorrectOptionID
        results = append(results, AnswerResult{
            QuestionID:    q.ID,
            SelectedOption: selectedOption,
            IsCorrect:     isCorrect,
        })
    }

    return results, nil
}
