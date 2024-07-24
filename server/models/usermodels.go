package models

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//root:Core@123@/blogweb?
// postgres://core:12345678@localhost:5432/cloud

var db *gorm.DB

func ConnectDB() error {
	dsn := "host=localhost user=core password=12345678 dbname=cloud"
	dbget, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db = dbget
	return nil

}

// DONE
func FindUser(username string) bool {
	var res int32
	db.Raw("SELECT COUNT(username) FROM users WHERE username=?", username).Scan(&res)
	return res != 0
}

func GetUsername(userid uint64) string {
	var u string
	db.Raw("SELECT username FROM users WHERE user_id=?", userid).Scan(&u)
	return u
}

// DONE
func CreateUser(user Users) (uint64, error) {
	var userid uint64
	tx := db.Begin()
	sql := "INSERT INTO users (username,email,password) VALUES(?,?,?,?,?) RETURNING user_id"
	r := tx.Raw(sql, user.Username, user.Email, user.Password).Scan(&userid)
	if r.Error != nil {
		tx.Rollback()
		return userid, r.Error
	} else {
		tx.Commit()
	}
	return userid, nil
}

func LoginUser(username string) (string, bool, uint64) {
	var exists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM users WHERE username=?)", username).Scan(&exists)
	if exists {
		var res Passanduserid
		db.Raw("SELECT user_id,password FROM users WHERE username= ?", username).Scan(&res)
		return res.Password, true, res.User_id
	}
	return "", false, 0
}

func UpdateUserAbout(about string, userid uint64) error {
	r := db.Raw("UPDATE users SET about = ? WHERE user_id = ?", about, userid)
	if r.Error != nil {
		return r.Error
	}
	return nil
}
func GetUserInfo(username string) (Users, error) {
	var user Users
	r := db.Raw("SELECT user_id,username,about,createdat FROM users WHERE username=?", username).Scan(&user)
	if r.RowsAffected == 0 {
		user.UserID = -1
		return user, errors.New("user not found")
	}
	if r.Error != nil {
		return user, r.Error
	}
	return user, nil
}
func DeleteUser(userid uint64) error {
	r := db.Raw("DELETE FROM comments WHERE user_id=?", userid)
	if r.Error != nil {
		return r.Error
	}
	r = db.Raw("DELETE FROM posts WHERE author_id=?", userid)
	if r.Error != nil {
		return r.Error
	}

	r = db.Raw("DELETE FROM users WHERE user_id=?", userid)
	if r.Error != nil {
		return r.Error
	}
	return nil
}
func CheckUserLoggedIn(userid uint64) bool {
	var active bool
	db.Raw("SELECT active FROM users WHERE user_id=?", userid).Scan(&active)
	return active
}
func AddProfilePictureStoreURL(userid uint64, imageURLString string) error {
	r := db.Exec("UPDATE users SET profilePicture=? WHERE user_id=?", imageURLString, userid)
	return r.Error
}

func LogOut(userid uint64) bool {
	p := make(chan int, 1)
	var ok bool
	go func() {
		tx := db.Begin()
		r := tx.Exec("UPDATE users SET active=? WHERE user_id=?", false, userid)
		if r.Error != nil {
			fmt.Println("error while logging out user")
			p <- 1
			ok = false
		}
		ok = true
		p <- 1
	}()
	<-p
	return ok
}
func UpdatePass(pass string, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("UPDATE users SET password=? WHERE userid=?", pass, userid)
	if r.Error != nil {
		tx.Rollback()
		fmt.Println(r.Error)
		return r.Error
	} else {
		tx.Commit()
		return nil
	}
}

func CheckUserExists(userid uint64) (bool, error) {
	var exists bool
	r := db.Raw("SELECT EXISTS (SELECT 1 FROM users WHERE user_id=?)", userid).Scan(&exists)
	if r.Error != nil {
		return exists, r.Error
	}
	return exists, nil
}
