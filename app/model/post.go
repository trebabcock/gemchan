package model

import "time"

type Post struct {
	Content    string
	Board      string
	ID         string
	Time       string
	LastUpdate time.Time
}

type Comment struct {
	Content string
	ReplyTo string
	ID      string
	Time    string
	PostID  string
}
