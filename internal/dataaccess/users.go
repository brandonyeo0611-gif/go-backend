package users

import (
	"database/sql"
	"errors"
	"log"

	"github.com/CVWO/sample-go-app/internal/database"
	"github.com/CVWO/sample-go-app/internal/models"
)

var ErrorUserNameTaken = errors.New("username already taken")
var ErrorShortPost = errors.New("too short post")
var ErrorShortComment = errors.New("too short comment")

func List(db *database.Database) ([]models.User, error) {
	users := []models.User{
		{
			UserID:   1,
			Username: "CVWO",
		},
	}
	return users, nil
}

func CreateUser(user models.User) error {
	db, err := database.GetDB()
	if err != nil {
		log.Println("CreateUser select error:", err)
		return err
	}

	var id int                                                                                      // create a variable of id first
	err = db.Conn.QueryRow("SELECT user_id FROM users WHERE username=$1 ", user.Username).Scan(&id) // fills id with existing user's user_id
	if err == nil {
		return ErrorUserNameTaken
	} else if err != sql.ErrNoRows {
		log.Println("CreateUser select error:", err)
		return err
	} // if there is an error and error is not errNoRows
	// no err signifys there is an existing user alr

	_, err = db.Conn.Exec("INSERT INTO users (username) VALUES ($1)", user.Username)
	if err != nil {
		log.Println("CreateUser select error:", err)
		return err
	}
	return nil
}

func UserLogin(user models.User) error {
	db, err := database.GetDB()
	if err != nil {
		log.Println("CreateUser select error:", err)
		return err
	}

	var id int
	err = db.Conn.QueryRow(
		`SELECT user_id
		FROM users
		WHERE username = $1`,
		user.Username,
	).Scan(&id)
	if err != nil {
		log.Println("CreateUser select error:", err)
		return err
	}
	return nil
}
// no. of arguments in scan same as no. of arguments in select

func CreatePost(post models.Post) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}
	if len(post.Content) < 100 {
		return ErrorShortPost
	}

	_, err = db.Conn.Exec("INSERT INTO post (user_id, content, content_type,title) VALUES ($1,$2,$3,$4)", post.UserID, post.Content, post.ContentType, post.Title)
	if err != nil {
		return err
	}
	return nil

}

func CreateComment(comment models.Comment) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}
	if len(comment.Content) < 50 {
		return ErrorShortComment
	}

	_, err = db.Conn.Exec("INSERT INTO comment (post_id,user_id, content) VALUES ($1,$2,$3)", comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return err
	}
	return nil

}

func UpdateLikesPost(postlikes models.PostLikes) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}

	var have int
	err = db.Conn.QueryRow(
		`SELECT like_value FROM post_likes
		WHERE post_id = $1 AND user_id = $2`,
		postlikes.PostID, postlikes.UserID,
	).Scan(&have)

	if err == sql.ErrNoRows {
		_, err = db.Conn.Exec(
			"INSERT INTO post_likes (post_id, user_id, like_value) VALUES ($1,$2,$3)", postlikes.PostID, postlikes.UserID, postlikes.LikeValue,
		)
		return err
	}

	if err != nil {
		return err
	}
	if have == postlikes.LikeValue {
		_, err = db.Conn.Exec(
			`DELETE FROM post_likes 
			WHERE post_id = $1 
			AND user_id = $2`, postlikes.PostID, postlikes.UserID,
		)
		return err
	}

	_, err = db.Conn.Exec(
		`UPDATE post_likes 
		SET like_value = $1 
		WHERE post_id = $2 
		AND user_id = $3`, postlikes.LikeValue, postlikes.PostID, postlikes.UserID,
	)
	return err

}

func UpdateLikesComment(commentlikes models.CommentLikes) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}
	var have int
	err = db.Conn.QueryRow(
		`SELECT like_value 
		FROM comment_likes 
		WHERE comment_id = $1 
		AND user_id = $2`, commentlikes.CommentID, commentlikes.UserID,
	).Scan(&have)

	if err == sql.ErrNoRows {
		_, err = db.Conn.Exec(
			`INSERT INTO comment_likes (comment_id, user_id, like_value)
			VALUES ($1,$2,$3)
			`, commentlikes.CommentID, commentlikes.UserID, commentlikes.LikeValue,
		)
		return err
	}

	if have == commentlikes.LikeValue {
		_, err = db.Conn.Exec(
			`DELETE FROM comment_likes 
			WHERE comment_id= $1 
			AND user_id = $2`, commentlikes.CommentID, commentlikes.UserID,
		)
		if err != nil {
			return err
		}
	}

	_, err = db.Conn.Exec(
		`UPDATE comment_likes
		SET like_value = $1
		WHERE comment_id = $2
		AND user_id = $3`, commentlikes.LikeValue, commentlikes.CommentID, commentlikes.UserID,
	)
	return err
}

func PostsByCategory(category string) ([]models.Post, error) {

	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Conn.Query(
		`SELECT post_id, user_id, content, created_at, content_type, title
		FROM post
		WHERE content_type = $1
		ORDER BY created_at DESC`,
		category,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// prevents data connection leakage
	// what if in middle of for loop then got error

	posts := []models.Post{}
	// make empty posts

	for rows.Next() {
		var p models.Post
		err = rows.Scan(
			&p.PostID,
			&p.UserID,
			&p.Content,
			&p.CreatedAt,
			&p.ContentType,
			&p.Title,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// checks for erros that happened during iterations, not when querying or scanning

	return posts, nil

}

func PostComments(postID string) ([]models.Comment, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Conn.Query(
		`SELECT comment_id, post_id, user_id, content, created_at
		FROM comment
		WHERE post_id = $1
		ORDER BY created_at DESC`, postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}

	for rows.Next() {
		var c models.Comment
		err = rows.Scan(
			&c.CommentID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, err
}

// SQL code
// $1 is temporary placeholder, sorta like format string
