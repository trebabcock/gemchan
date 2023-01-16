package handler

import (
	"gemchan/app/model"
	"log"
	"sort"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

func CreatePost(db *gorm.DB, content, board string) string {
	id := shortuuid.New()
	post := model.Post{
		Content:    content,
		Board:      board,
		ID:         id,
		Time:       time.Now().UTC().Format(time.RFC1123),
		LastUpdate: time.Now().UTC(),
	}
	if err := db.Create(&post).Error; err != nil {
		log.Println("Couldn't save post:", err)
		return "Unable to save post"
	}
	return id
}

func GetAllPosts(db *gorm.DB) []model.Post {
	posts := []model.Post{}
	db.Find(&posts)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].LastUpdate.After(posts[j].LastUpdate)
	})
	return posts
}

func GetPostsFromBoard(db *gorm.DB, board string) []model.Post {
	posts := []model.Post{}
	db.Find(&posts, db.Where(model.Post{Board: board}))
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].LastUpdate.After(posts[j].LastUpdate)
	})
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
	id := shortuuid.New()
	c := model.Comment{
		Content: content,
		ReplyTo: "",
		ID:      id,
		Time:    time.Now().UTC().Format(time.RFC1123),
		PostID:  postID,
	}
	p.LastUpdate = time.Now().UTC()
	if err := db.Save(&p).Error; err != nil {
		log.Println("Couldn't add comment:", err)
		return
	}
	if err := db.Save(&c).Error; err != nil {
		log.Println("Couldn't add comment:", err)
		return
	}
}

func AddCommentReply(db *gorm.DB, content, replyto, postID string) {
	p := GetPost(db, postID)
	id := shortuuid.New()
	c := model.Comment{
		Content: content,
		ReplyTo: replyto,
		ID:      id,
		Time:    time.Now().UTC().Format(time.RFC1123),
		PostID:  postID,
	}

	p.LastUpdate = time.Now().UTC()
	if err := db.Save(&p).Error; err != nil {
		log.Println("Couldn't add comment:", err)
		return
	}
	if err := db.Save(&c).Error; err != nil {
		log.Println("Couldn't add comment:", err)
		return
	}
}

func GetComments(db *gorm.DB, postID string) []model.Comment {
	c := []model.Comment{}
	db.Find(&c, db.Where(model.Comment{PostID: postID}))
	return c
}
