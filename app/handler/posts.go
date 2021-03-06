package handler

import (
	"gem/app/model"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePost(db *gorm.DB, content, board string) string {
	id := uuid.NewString()
	post := model.Post{
		Content:    content,
		Board:      board,
		ID:         id,
		Time:       time.Now().UTC().Format(time.Stamp),
		LastUpdate: time.Now().UTC(),
	}
	if err := db.Create(&post).Error; err != nil {
		log.Println("Couldn't save post:", err)
		return "Unable to save post"
	}
	return id
}

func GetPostsFromBoard(db *gorm.DB, board string) []model.Post {
	posts := []model.Post{}
	db.Find(&posts, db.Where(model.Post{Board: board}))
	return posts
}

func GetPost(db *gorm.DB, id string) *model.Post {
	post := model.Post{}
	if err := db.First(&post, db.Where(model.Post{ID: id})).Error; err != nil {
		log.Println("Couldn't find post:", err)
		return nil
	}
	return &post
}

func AddComment(db *gorm.DB, content, postID string) {
	p := GetPost(db, postID)
	id := uuid.NewString()
	c := model.Comment{
		Content: content,
		ReplyTo: "",
		ID:      id,
		Time:    time.Now().UTC().Format(time.Stamp),
		PostID:  postID,
	}
	p.Comments = append(p.Comments, c)
	p.LastUpdate = time.Now().UTC()
	if err := db.Save(&p).Error; err != nil {
		log.Println("Couldn't add comment:", err)
		return
	}
}

func AddCommentReply(db *gorm.DB, content, replyto, postID string) {
	p := GetPost(db, postID)
	id := uuid.NewString()
	c := model.Comment{
		Content: content,
		ReplyTo: replyto,
		ID:      id,
		Time:    time.Now().UTC().Format(time.Stamp),
		PostID:  postID,
	}
	p.Comments = append(p.Comments, c)
	p.LastUpdate = time.Now().UTC()
	if err := db.Save(&p).Error; err != nil {
		log.Println("Couldn't add comment:", err)
		return
	}
}

func GetComments(db *gorm.DB, postID string) []model.Comment {
	c := []model.Comment{}
	db.Find(&c, db.Where(model.Comment{PostID: postID}))
	return c
}
