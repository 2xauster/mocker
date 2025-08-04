package entities

import "time"

type Session struct {
	ID string `type:"TEXT" cnstr:"PRIMARY KEY"`

	UserID string `type:"TEXT" cnstr:"NOT NULL" ref:"User(ID)"`
	MockID string `type:"TEXT" cnstr:"NOT NULL" ref:"Mock(ID)"`
	TTL    int    `type:"NUMBER" cnstr:"NOT NULL" ref:"Mock(TimeMins)"`

	CreatedAt     time.Time `type:"TEXT" cnstr:"NOT NULL"`
}