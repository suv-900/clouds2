package models

import (
	"errors"
	"fmt"
)

func CreatePost(post Posts) (uint64, error) {

	userexists := CheckUserExists(post.Author_id)
	if !userexists {
		fmt.Println("User doesnt exists.Post Creation Failed")
		return 0, errors.New("user doesnt exists.Post Creation Failed")
	}

	var postid uint64
	tx := db.Begin()
	r := tx.Raw("INSERT INTO posts (post_title,post_content,author_id,post_likes) VALUES(?,?,?,?) RETURNING post_id", post.Post_title, post.Post_content, post.Author_id, 0).Scan(&postid)
	if r.Error != nil {
		tx.Rollback()
		return postid, r.Error
	} else {
		tx.Commit()
	}
	return postid, nil
}

// deletes a post
func DeletePost(postid uint64) error {

	r := db.Exec("DELETE FROM posts WHERE postid=?", postid)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func GetAllPosts() []Posts {
	var posts []Posts
	db.Raw("SELECT * FROM posts").Scan(&posts)
	return posts
}

func UpdatePost(postid uint64, post Posts) error {

	tx := db.Begin()
	r := tx.Exec("UPDATE posts SET post_content=? post_title=? WHERE postid=?", post.Post_content, post.Post_title, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
		return nil
	}

}
func GetPostsByUserId(userid uint64) []Posts {
	var posts []Posts
	db.Raw("SELECT (post_title,post_content) FROM posts WHERE authorid=? LIMIT 5", userid).Scan(&posts)
	return posts
}

func GetPosts(offset uint64) []Posts {
	var posts []Posts
	db.Raw(`SELECT * FROM posts ORDER BY post_likes DESC LIMIT 5 OFFSET ?`, offset).Scan(&posts)
	return posts
}

// func GetPostAndUserPreferences(postid uint64, userid uint64) (Posts, error) {
// 	var post Posts
// 	r := db.Raw("SELECT * FROM posts WHERE post_id=?", postid).Scan(&post)
// 	if r.Error != nil {
// 		return post,r.Error
// 	}
// 	return post, nil
// }

func CheckUserReaction(userid uint64, postid uint64) (bool, bool, error) {
	var userLikedPost bool
	var userDislikedPost bool
	var count1 int
	var count2 int
	var err error
	r := db.Raw("SELECT COUNT(*) from posts_liked_by_users WHERE user_id=? AND post_id=?", userid, postid).Scan(&count1)
	if r.Error != nil {
		return userLikedPost, userDislikedPost, r.Error
	}
	r = db.Raw("SELECT COUNT(*) FROM posts_disliked_by_users WHERE user_id=? AND post_id=?", userid, postid).Scan(&count2)
	if r.Error != nil {
		return userLikedPost, userDislikedPost, r.Error
	}

	if count1 == 1 && count2 == 2 {
		return userLikedPost, userDislikedPost, errors.New("illegalState:User reaction like and dislike both exists.check CheckUserReaction() function")
	}

	if count1 == 1 {
		userLikedPost = true
		userDislikedPost = false
	} else if count2 == 1 {
		userLikedPost = false
		userDislikedPost = true
	}
	return userLikedPost, userDislikedPost, err
}
func PostById(postid uint64) (Posts, error) {
	var post Posts
	r := db.Raw(`SELECT * FROM posts WHERE post_id=?`, postid).Scan(&post)
	if r.Error != nil {
		return post, r.Error
	}
	return post, r.Error
}

// server should be dumb and do what said nothing extra
// func LikePostByID(userid uint64, postid uint64) error {
// 	var err error

// 	c := make(chan int, 1)
// 	go func() {
// 		r := db.Exec("DELETE FROM posts_disliked_by_user WHERE user_id=? AND post_id=?", userid, postid)
// 		if r.Error != nil {
// 			err = r.Error
// 			c <- 1
// 			return
// 		}
// 		c <- 1
// 	}()
// 	<-c

// 	if err != nil {
// 		return err
// 	}

// 	a := make(chan int, 1)
// 	go func() {
// 		tx := db.Begin()
// 		r := tx.Exec("UPDATE posts SET post_likes=post_likes+1 WHERE post_id=?", postid)
// 		if r.Error != nil {
// 			err = r.Error
// 			a <- 1
// 			tx.Rollback()
// 			return
// 		}
// 		tx.Commit()
// 		a <- 1
// 	}()
// 	<-a
// 	if err != nil {
// 		return err
// 	}

// 	b := make(chan int, 1)
// 	go func() {
// 		tx := db.Begin()
// 		r := tx.Exec("INSERT INTO posts_liked_by_user (user_id,post_id,liked) VALUES(?,?,?)", userid, postid, true)
// 		if r.Error != nil {
// 			err = r.Error
// 			b <- 1
// 			tx.Rollback()
// 			return
// 		}
// 		tx.Commit()
// 		b <- 1
// 	}()
// 	<-b
// 	return err
// }

func LikePost(postid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("INSERT INTO posts_liked_by_users(user_id,post_id) VALUES(?,?) ", userid, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE posts SET post_likes= post_likes + 1 WHERE post_id=?", postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}

func DislikePost(postid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("INSERT INTO posts_disliked_by_users(user_id,post_id) VALUES(?,?) ", userid, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE posts SET post_likes=post_likes - 1 WHERE post_id=?", postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}
func RemoveLikeFromPost(postid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("DELETE FROM posts_liked_by_users WHERE user_id=? AND post_id=?", userid, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE posts SET post_likes=post_likes - 1 WHERE post_id=?", postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}

func RemoveDislikeFromPost(postid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("DELETE FROM posts_disliked_by_users WHERE user_id=? AND post_id=?", userid, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE posts SET post_likes=post_likes + 1 WHERE post_id=?", postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}

func FeedGenerator(userid uint64) []Posts {

	var categories []string

	db.Raw(`SELECT c.category_name
  FROM post_likes pl JOIN posts p ON pl.post_id=p.post_id
  JOIN categories c ON c.category_id=p.category_id
  WHERE pl.user_id=?
  GROUP BY c.category_name
  LIMIT 5`, userid).Scan(&categories)

	var posts []Posts
	db.Raw(` 
  SELECT p.author_id,p.post_title,p.post_content,COUNT(pl.like_id) as likes_num
  FROM post_likes pl JOIN post p ON pl.post_id=p.post_id
  JOIN categories c ON c.category_id=p.category_id
  WHERE c.category_name=?
  GROUP BY p.post_id,p.author_id,p.post_title,p.post_content
  ORDER BY likes_num DESC
  LIMIT 10
  `, categories).Scan(&posts)
	return posts
}
