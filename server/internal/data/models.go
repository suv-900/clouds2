package data

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	ErrConflict            = errors.New("already exists")
	ErrRecordNotFound      = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")
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
}

func GetModels(gd *gorm.DB, mc *mongo.Client) Models {
	return Models{
		Users: UserClient{db: gd},
	}
}
