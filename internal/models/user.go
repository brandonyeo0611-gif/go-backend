package models
import "time"
type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

type Comment struct {
	CommentID string `json:"comment_id"`
	PostID string `json:"post_id"`
	UserID int `json:"user_id"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	PostID string `json:"post_id"`
	UserID int `json:"user_id"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ContentType string `json:"content_type"`
	Title string `json:"title"`
}

type PostLikes struct {
	PostID string `json:"post_id"`
	UserID int `json:"user_id"`
	LikeValue int `json:"like_value"`
}

type CommentLikes struct {
	CommentID string `json:"comment_id"`
	UserID int `json:"user_id"`
	LikeValue int `json:"like_value"`
}

// json tag match table col names 
// struct are exported so need capitalise