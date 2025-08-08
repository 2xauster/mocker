package session

import (
	"database/sql"

	"github.com/ashtonx86/mocker/internal/data"
)

type SessionManager struct {
	ActiveSessions map[string]Session // [K : Session.ID] : [v : Session]
	ExpiredSessions []Session // sessions to get rid of 

	DB *sql.DB
	Redis *data.Redis
}

func NewSessionManager(db *sql.DB, redisClient *data.Redis) *SessionManager {
	return &SessionManager{
		ActiveSessions: make(map[string]Session),
		ExpiredSessions: []Session{},

		DB: db,
		Redis: redisClient,
	}
}

// Create new session
func (s *SessionManager) New() {
	
}