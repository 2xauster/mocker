package session

import "time"

type Session struct {
	ID     string `json:"id"`
	MockID string `json:"mock_id"`
	UserID string `json:"user_id"`

	TTL       int `json:"ttl"`
	Answers   map[string]string // [K : questionID] [V : optionID/answerID]
	
	CreatedAt time.Time `json:"created_at"`
}