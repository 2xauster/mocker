package entities

import "time"

// Represent the "mock" table.
type Mock struct {
	ID string `type:"TEXT" cnstr:"PRIMARY KEY"`
	Topic string `type:"TEXT" cnstr:"NOT NULL"`
	Instructions string `type:"TEXT"`
	TimeMins int `type:"NUMBER" cnstr:"NOT NULL"`
	
	AuthorID string `type:"TEXT" cnstr:"NOT NULL" ref:"User(ID)"`

	CreatedAt time.Time `type:"TEXT" cnstr:"NOT NULL"`
	LastUpdatedAt time.Time `type:"TEXT" cnstr:"NOT NULL"`
}

type MockQuestion struct {
	ID string `type:"TEXT" cnstr:"PRIMARY KEY"`
	Problem string `type:"TEXT" cnstr:"NOT NULL"`
	Points int `type:"NUMBER" cnstr:"NOT NULL"`
	
	CorrectOptionID string `type:"TEXT" cnstr:"NOT NULL" ref:"MockOption(ID)"`
	MockID string `type:"TEXT" cnstr:"NOT NULL" ref:"Mock(ID)"`

	CreatedAt time.Time `type:"TEXT" cnstr:"NOT NULL"`
	LastUpdatedAt time.Time `type:"TEXT" cnstr:"NOT NULL"` 
}

type MockOption struct {
	ID string `type:"TEXT" cnstr:"PRIMARY KEY"`
	Number int `type:"NUMBER" cnstr:"NOT NULL"`
	Option string `type:"TEXT" cnstr:"NOT NULL"`

	QuestionID string `type:"TEXT" cnstr:"NOT NULL" ref:"MockQuestion(ID)"`

	CreatedAt time.Time `type:"TEXT" cnstr:"NOT NULL"`
	LastUpdatedAt time.Time `type:"TEXT" cnstr:"NOT NULL"`
}