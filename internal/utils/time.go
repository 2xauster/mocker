package utils

import (
	"fmt"
	"time"
)

const Layout = "2006-01-02 15:04:05.9999999-07:00"

func ParseTime(timeString string) (*time.Time, error){
	t, err := time.Parse(Layout, timeString)
	if err != nil {
		 return nil, fmt.Errorf("invalid createdAt format: %w", err)
	}
	return &t, err 
}