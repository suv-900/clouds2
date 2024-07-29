package data

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	ErrConflict            = errors.New("already exists")
	ErrRecordNotFound      = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")
	ErrIllegalState        = errors.New("illegal state")
)

const context_timeout = 5 * time.Second

type Models struct {
	Users interface {
		AddUser(cx context.Context, user *User) error
		GetUser(cx context.Context, username string) (User, error)
		UpdateUser(cx context.Context, user *User) error
		DeleteUser(cx context.Context, userID uint64) error

		CheckUserExists(cx context.Context, username string) (bool, error)
		GetUserPassword(cx context.Context, username string) (string, error)
	}
	Posts interface {
		AddPost(cx context.Context, post *Post) (int, error)
		GetPost(cx context.Context, postID primitive.ObjectID) (Post, error)
		UpdatePost(cx context.Context, update bson.D, postID primitive.ObjectID) error
		DeletePost(cx context.Context, postID primitive.ObjectID) error

		CheckPostTitleExists(cx context.Context, post_title string) (bool, error)
		GetPostsByAuthorID(cx context.Context, authorID int64, offset int64) ([]Post, error)
		IncrementLike(cx context.Context, postID primitive.ObjectID) error
		DecrementLike(cx context.Context, postID primitive.ObjectID) error
	}
}

func GetModels(gd *gorm.DB, mc *mongo.Client) Models {
	return Models{
		Users: UserClient{db: gd},
		Posts: PostClient{client: mc},
	}
}
