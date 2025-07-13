package entities

import "time"

type User struct {
	ID           string `type:"TEXT" cnstr:"PRIMARY KEY"`
	Name         string `type:"VARCHAR(45)" cnstr:"NOT NULL"`
	Email        string `type:"TEXT" cnstr:"UNIQUE NOT NULL"`
	PasswordHash string `type:"TEXT" cnstr:"NOT NULL"`
	
	CreatedAt time.Time `type:"TEXT" cnstr:"NOT NULL"`
	LastUpdatedAt time.Time `type:"TEXT" cnstr:"NOT NULL"`
}