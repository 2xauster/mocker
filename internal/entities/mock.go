package entities

import "time"

// Represent the "mock" table.
type Mock struct {
	ID             string    `type:"TEXT" cnstr:"PRIMARY KEY" json:"id"`
	Topic          string    `type:"TEXT" cnstr:"NOT NULL" json:"topic"`
	Instructions   string    `type:"TEXT" json:"instructions"`
	TimeMins       int       `type:"NUMBER" cnstr:"NOT NULL" json:"time_mins"`
	AuthorID       string    `type:"TEXT" cnstr:"NOT NULL" ref:"User(ID)" json:"author_id"`
	CreatedAt      time.Time `type:"TEXT" cnstr:"NOT NULL" json:"created_at"`
	LastUpdatedAt  time.Time `type:"TEXT" cnstr:"NOT NULL" json:"last_updated_at"`
}

type MockQuestion struct {
	ID              string    `type:"TEXT" cnstr:"PRIMARY KEY" json:"id"`
	Problem         string    `type:"TEXT" cnstr:"NOT NULL" json:"problem"`
	Points          int       `type:"NUMBER" cnstr:"NOT NULL" json:"points"`
	CorrectOptionID string    `type:"TEXT" cnstr:"NOT NULL" ref:"MockOption(ID)" json:"correct_option_id,omitempty"`
	MockID          string    `type:"TEXT" cnstr:"NOT NULL" ref:"Mock(ID)" json:"mock_id"`
	CreatedAt       time.Time `type:"TEXT" cnstr:"NOT NULL" json:"created_at"`
	LastUpdatedAt   time.Time `type:"TEXT" cnstr:"NOT NULL" json:"last_updated_at"`
}

type MockOption struct {
	ID            string    `type:"TEXT" cnstr:"PRIMARY KEY" json:"id"`
	Number        int       `type:"NUMBER" cnstr:"NOT NULL" json:"number"`
	Option        string    `type:"TEXT" cnstr:"NOT NULL" json:"option"`
	QuestionID    string    `type:"TEXT" cnstr:"NOT NULL" ref:"MockQuestion(ID)" json:"question_id"`
	CreatedAt     time.Time `type:"TEXT" cnstr:"NOT NULL" json:"created_at"`
	LastUpdatedAt time.Time `type:"TEXT" cnstr:"NOT NULL" json:"last_updated_at"`
}