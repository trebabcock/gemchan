package model

import "time"

type Post struct {
	Content    string
	Board      string
	ID         string
	Time       string
	LastUpdate time.Time
	Comments   []Comment
}

type Comment struct {
	Content string
	ReplyTo string
	ID      string
	Time    string
	PostID  string
}
