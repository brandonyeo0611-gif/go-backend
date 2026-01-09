package users

import (
	"database/sql"
	"errors"
	"fmt"
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

func CreateUser(user models.User, db *database.Database) error {

	var id int                                                                                       // create a variable of id first
	err := db.Conn.QueryRow("SELECT user_id FROM users WHERE username=$1 ", user.Username).Scan(&id) // fills id with existing user's user_id
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

func UserLogin(user models.User, db *database.Database) (int, error) {

	var id int
	err := db.Conn.QueryRow(
		`SELECT user_id
		FROM users
		WHERE username = $1`,
		user.Username,
	).Scan(&id)
	if err != nil {
		log.Println("CreateUser select error:", err)
		return 0, err
	}
	return id, nil
}

func ChangeProfilePic(user models.User, db *database.Database) error {
	_, err := db.Conn.Exec(
		`UPDATE users
		SET profile_url = $1
		WHERE username = $2 AND user_id = $3`,
		*user.ProfileURL, user.Username, user.UserID,
	)
	if err != nil {
		return err
	}
	return nil
}

// no. of arguments in scan same as no. of arguments in select

func CreatePost(post models.Post, db *database.Database) error {
	if len(post.Content) < 100 {
		return ErrorShortPost
	}

	_, err := db.Conn.Exec("INSERT INTO post (user_id, username, content, content_type,title) VALUES ($1,$2,$3,$4,$5)", post.UserID, post.Username, post.Content, post.ContentType, post.Title)
	if err != nil {
		return err
	}
	return nil

}

func CreateComment(comment models.Comment, db *database.Database) error {
	if len(comment.Content) < 50 {
		return ErrorShortComment
	}

	_, err := db.Conn.Exec("INSERT INTO comment (post_id,user_id, content) VALUES ($1,$2,$3)", comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return err
	}
	return nil

}

func UpdateLikesPost(postlikes models.PostLikes, db *database.Database) error {

	var have int
	err := db.Conn.QueryRow(
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

func UpdateLikesComment(commentlikes models.CommentLikes, db *database.Database) error {

	var have int
	err := db.Conn.QueryRow(
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

func PostsByCategory(category string, relevance string, db *database.Database) ([]models.FullPostStruct, error) {
	orderBy := "p.created_at DESC" // default

	switch relevance {
	case "leastRecent":
		orderBy = "p.created_at ASC"
	case "mostLikes":
		orderBy = "likes DESC"
	case "leastLikes":
		orderBy = "likes ASC"
	}
	if category == "All" {

		query := fmt.Sprintf(`
SELECT p.post_id, p.user_id, p.username, p.content, p.created_at, p.content_type, p.title,
COALESCE(SUM(p1.like_value), 0) AS likes
FROM post p
LEFT JOIN post_likes p1 ON p.post_id = p1.post_id
WHERE content_type = $1
GROUP BY p.post_id, p.user_id, p.username, p.content, p.created_at, p.content_type, p.title
ORDER BY %s`, orderBy) // learning to inject with sprintf (values use placeholder $1,$2,$3 but ORDER BY or column names can use sprintf)

		rows, err := db.Conn.Query(query)

		if err != nil {
			return nil, err
		}
		defer rows.Close()
		// prevents data connection leakage
		// what if in middle of for loop then got error

		posts := []models.FullPostStruct{}
		// make empty posts

		for rows.Next() {
			var p models.FullPostStruct
			err = rows.Scan(
				&p.PostID,
				&p.UserID,
				&p.Username,
				&p.Content,
				&p.CreatedAt,
				&p.ContentType,
				&p.Title,
				&p.Likes,
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

	query := fmt.Sprintf(`
SELECT p.post_id, p.user_id, p.username, p.content, p.created_at, p.content_type, p.title,
COALESCE(SUM(p1.like_value), 0) AS likes
FROM post p
LEFT JOIN post_likes p1 ON p.post_id = p1.post_id
WHERE content_type = $1
GROUP BY p.post_id, p.user_id, p.username, p.content, p.created_at, p.content_type, p.title
ORDER BY %s`, orderBy)

	rows, err := db.Conn.Query(query, category)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// prevents data connection leakage
	// what if in middle of for loop then got error

	posts := []models.FullPostStruct{}
	// make empty posts

	for rows.Next() {
		var p models.FullPostStruct
		err = rows.Scan(
			&p.PostID,
			&p.UserID,
			&p.Username,
			&p.Content,
			&p.CreatedAt,
			&p.ContentType,
			&p.Title,
			&p.Likes,
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

func PostComments(postID string, db *database.Database) ([]models.CommentResponse, error) {

	rows, err := db.Conn.Query(
		`SELECT c.comment_id, c.post_id, c.user_id, c.content, c.created_at, u.username
		FROM comment c
		INNER JOIN users u ON c.user_id = u.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC`, postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.CommentResponse{}

	for rows.Next() {
		var c models.CommentResponse
		err = rows.Scan(
			&c.CommentID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.Username,
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

func GetPost(postID string, db *database.Database) (*models.Post, error) {

	var post models.Post

	err := db.Conn.QueryRow(
		`SELECT post_id, user_id, username, content, created_at, content_type, title
		FROM post 
		WHERE post_id = $1`,
		postID,
	).Scan(&post.PostID,
		&post.UserID,
		&post.Username,
		&post.Content,
		&post.CreatedAt,
		&post.ContentType,
		&post.Title,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func GetIndividualLike(userID int, postID string, db *database.Database) (int, error) {

	var likeValue int
	err := db.Conn.QueryRow(
		`SELECT like_value FROM post_likes WHERE post_id = $1 AND user_id = $2`, postID, userID,
	).Scan(&likeValue)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err // gpt say this one clearer though no diff
	}
	return likeValue, err
}

func RefreshAccessToken(RefreshToken string) {

}

// SQL code
// $1 is temporary placeholder, sorta like format string
