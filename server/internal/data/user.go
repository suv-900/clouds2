package models

import (
	"context"
	"errors"
	"time"

	"gorm.io/datatypes"

	"gorm.io/gorm"
)

// postgres://core:12345678@localhost:5432/cloud
type User struct {
	ID uint `gorm:"primaryKey"`

	Username      string
	Active        bool
	Email         string
	Password      string
	BirthDate     datatypes.Date
	ProfilePicURL string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserClient struct {
	db *gorm.DB
}

var (
	ErrConflict            = errors.New("already exists")
	ErrRecordNotFound      = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")
)

func (u UserClient) CheckUserExists(c context.Context, username string) (bool, error) {
	cx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	sql := `select exists (select 1 from users where username = ?) as exists`
	var found bool

	t := u.db.WithContext(cx).Raw(sql, username).Scan(&found)
	if t.Error != nil {
		return false, t.Error
	}

	return found, nil
}

func (u UserClient) AddUser(c context.Context, user *User) error {
	cx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	t := u.db.WithContext(cx).Create(user)
	if t.Error != nil {
		if errors.Is(t.Error, gorm.ErrDuplicatedKey) {
			return ErrConflict
		}
		return ErrInternalServerError
	}

	return nil
}

func (u UserClient) GetUsername(c context.Context, userid uint64) (string, error) {
	cx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var username string
	sql := `select username from users where user_id = ?`
	t := u.db.WithContext(cx).Raw(sql, userid).Scan(&username)

	if t.Error != nil {
		if errors.Is(t.Error, gorm.ErrRecordNotFound) {
			return username, ErrRecordNotFound
		}
		return username, ErrInternalServerError
	}

	return username, nil
}

func (u UserClient) GetUserPassword(c context.Context, username string) (string, error) {
	cx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var dbpassword string
	sql := `SELECT password FROM users WHERE username = ?`
	t := u.db.WithContext(cx).Raw(sql, username).Scan(&dbpassword)

	if t.Error != nil {
		if errors.Is(t.Error, gorm.ErrRecordNotFound) {
			return dbpassword, ErrRecordNotFound
		}
		return dbpassword, ErrInternalServerError
	}

	return dbpassword, nil
}

func (u UserClient) GetUser(c context.Context, username string) (User, error) {
	cx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var user User
	sql := `SELECT user_id,about,active,email,createdat FROM users WHERE username = ?`

	t := u.db.WithContext(cx).Raw(sql, username).Scan(&user)

	if t.Error != nil {
		if errors.Is(t.Error, gorm.ErrRecordNotFound) {
			return user, ErrRecordNotFound
		}
		return user, ErrInternalServerError
	}
	return user, nil
}
func (u UserClient) UpdateProfilePictureURL(c context.Context, user *User, iurl string) error {
	cx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	t := u.db.WithContext(cx).Model(user).Update("ProfilePicURL", iurl)

	if t.Error != nil {
		if errors.Is(t.Error, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}

		return ErrInternalServerError
	}

	return nil
}

func (u UserClient) UpdateUserPassword(c context.Context, password string, user *User) error {
	cx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	t := u.db.WithContext(cx).Model(user).Update("password", password)

	if t.Error != nil {
		if errors.Is(t.Error, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return ErrInternalServerError
	}

	return nil
}

func (u UserClient) DeleteUser(c context.Context, userid uint64) error {
	cx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	t := u.db.WithContext(cx).Delete(&User{}, userid)

	if t.Error != nil {
		if errors.Is(t.Error, gorm.ErrRecordNotFound) {
			return ErrRecordNotFound
		}
		return ErrInternalServerError
	}

	return nil
}
