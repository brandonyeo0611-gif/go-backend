package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CVWO/sample-go-app/internal/api"
	users "github.com/CVWO/sample-go-app/internal/dataaccess"
	"github.com/CVWO/sample-go-app/internal/database"
	"github.com/CVWO/sample-go-app/internal/models"
	"github.com/pkg/errors"
)

const (
	ListUsers = "users.HandleList"

	SuccessfulListUsersMessage = "Successfully listed users"
	ErrRetrieveDatabase        = "Failed to retrieve database in %s"
	ErrRetrieveUsers           = "Failed to retrieve users in %s"
	ErrEncodeView              = "Failed to retrieve users in %s"
)

func HandleList(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	db, err := database.GetDB()

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveDatabase, ListUsers))
	}

	users, err := users.List(db)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveUsers, ListUsers))
	}

	data, err := json.Marshal(users)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, ListUsers))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulListUsersMessage},
	}, nil
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Fail to create user"},
			ErrorCode: 100,
		}, nil
	}

	if newUser.Username == "" {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Username cannot be empty"},
			ErrorCode: 101,
		}, nil
	}

	err = users.CreateUser(newUser)
	if err != nil {
		if err == users.ErrorUserNameTaken {
			return &api.Response{
				Payload:   api.Payload{},
				Messages:  []string{"Username taken"},
				ErrorCode: 102,
			}, nil
		} else {
			return &api.Response{
				Payload:   api.Payload{},
				Messages:  []string{"Fail to create user"},
				ErrorCode: 103,
			}, nil
		}
	}

	data, err := json.Marshal(newUser)
	if err != nil {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Fail to create user"},
			ErrorCode: 104,
		}, nil
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages:  []string{"User created successfully"},
		ErrorCode: 0,
	}, nil
}

func HandleUserLogin(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	var user models.User 
	err := json.NewDecoder(r.Body).Decode(&user) // decontruct the body to get the go code
	if err != nil {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Fail to login"},
			ErrorCode: 5000,
		}, nil
	} 

	err = users.UserLogin(user)
	if err != nil {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"wrong username"},
			ErrorCode: 1000,
		}, nil
	}
	return &api.Response{
		Payload: api.Payload{},
		Messages:  []string{"Correct username"},
		ErrorCode: 0,
	}, nil

}

func HandleCreatePost(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	var newPost models.Post
	err := json.NewDecoder(r.Body).Decode(&newPost) // decontruct the body to get the go code
	if err != nil {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Fail to create post"},
			ErrorCode: 5000,
		}, nil
	} 
	err = users.CreatePost(newPost)
	if err != nil {
		if errors.Is(err, users.ErrorShortPost) {
			return &api.Response{
				Payload:   api.Payload{},
				Messages:  []string{"Cannot create too short post"},
				ErrorCode: 1000,
			}, nil
		}
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Fail to create post"},
			ErrorCode: 5000,
		}, nil
	}
	return &api.Response{
		Payload:   api.Payload{},
		Messages:  []string{"Successfully create Post"},
		ErrorCode: 0,
	}, nil

}

func HandleCreateComment(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	var newComment models.Comment
	err := json.NewDecoder(r.Body).Decode(&newComment)
	if err != nil {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Fail to create post"},
			ErrorCode: 5000,
		}, nil
	}

	err = users.CreateComment(newComment)
	if err != nil {
		if errors.Is(err, users.ErrorShortComment) {
			return &api.Response{
				Payload:   api.Payload{},
				Messages:  []string{"Cannot create too short Comment"},
				ErrorCode: 1000,
			}, nil
		}
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Failed to create comment"},
			ErrorCode: 5000,
		}, nil
	}

	return &api.Response{
		Payload:   api.Payload{},
		Messages:  []string{"Successfully create Comment"},
		ErrorCode: 0,
	}, nil
}

func HandleLikesPost(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	var newLike models.PostLikes
	err := json.NewDecoder(r.Body).Decode(&newLike)
	if err != nil {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Fail to update likes"},
			ErrorCode: 5000,
		}, nil
	}
	err = users.UpdateLikesPost(newLike)
	if err != nil {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Fail to update likes"},
			ErrorCode: 5000,
		}, nil
	}

	return &api.Response{
		Payload:   api.Payload{},
		Messages:  []string{"Successfully update likes"},
		ErrorCode: 0,
	}, nil
}

func HandleLikesComment(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	var newLike models.CommentLikes
	err := json.NewDecoder(r.Body).Decode(&newLike)
	if err != nil {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Fail to update likes"},
			ErrorCode: 5000,
		}, nil
	}

	err = users.UpdateLikesComment(newLike)
	if err != nil {
		return &api.Response{
			Payload:   api.Payload{},
			Messages:  []string{"Fail to update likes"},
			ErrorCode: 5000,
		}, nil
	}
	return &api.Response{
		Payload:   api.Payload{},
		Messages:  []string{"Successfully update likes"},
		ErrorCode: 0,
	}, nil
}

func HandleGetPostsByCategory(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	category := r.URL.Query().Get("category")
	posts, err := users.PostsByCategory(category) 
	if err != nil {
		return &api.Response{
            Payload:   api.Payload{},
            Messages:  []string{"Failed to fetch posts"},
            ErrorCode: 5000,
        }, nil
	}
	data, err := json.Marshal(posts)
	if err != nil {
		return &api.Response{
            Payload:   api.Payload{},
            Messages:  []string{"Failed to fetch posts"},
            ErrorCode: 5001,
        }, nil
	}
	return &api.Response{
            Payload:   api.Payload{Data: data,},
            Messages:  []string{"Successfully fetch posts"},
            ErrorCode: 0,
        }, nil
}

func HandleGetComment(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	postID := r.URL.Query().Get("post")
	comments, err := users.PostComments(postID)
	if err != nil {
		return &api.Response{
            Payload:   api.Payload{},
            Messages:  []string{"Failed to fetch comments"},
            ErrorCode: 5000,
        }, nil
	}

	data, err := json.Marshal(comments)
	if err != nil {
		return &api.Response{
            Payload:   api.Payload{},
            Messages:  []string{"Failed to fetch comments"},
            ErrorCode: 5001,
        }, nil
	}
	return &api.Response{
            Payload:   api.Payload{Data: data,},
            Messages:  []string{"Successfully fetch comments"},
            ErrorCode: 0,
        }, nil
}

// remember that handlers does not decide the logic.. it just see if got error or not, should remain as short as possible
// the data access handles the logic
