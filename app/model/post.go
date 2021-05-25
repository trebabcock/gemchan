package model

import "time"

type Post struct {
	Content    string
	Board      string
	ID         string
	Time       string
	LastUpdate time.Time
	Comments   []Comment `gorm:"foreignKey:ID;references:ID"`
}

type Comment struct {
	Content string
	ID      string
	Time    string
	PostID  string
}
