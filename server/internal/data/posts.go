package models

import (
	"errors"
	"fmt"
)

func CreatePost(post Posts) (uint64, error) {

	var postid uint64
	tx := db.Begin()
	r := tx.Raw("INSERT INTO posts (post_title,post_content,author_id,post_likes,author_name) VALUES(?,?,?,?,?) RETURNING post_id", post.Post_title, post.Post_content, post.Author_id, 0, post.Author_name).Scan(&postid)
	if r.Error != nil {
		tx.Rollback()
		return postid, r.Error
	} else {
		tx.Commit()
	}
	return postid, nil
}

// deletes a post
func DeletePosts(posts []Posts) error {
	fmt.Println("posts", posts)
	r := db.Delete(&posts)
	if r.Error != nil {
		return r.Error
	}
	return nil
}
func GetFeaturedPosts(offset uint64) []Posts {
	var postMetaData []Posts
	limit := 3
	sql := "SELECT post_id,post_title FROM posts ORDER BY post_likes DESC OFFSET ? LIMIT ?"
	db.Raw(sql, offset, limit).Scan(&postMetaData)
	return postMetaData
}
func GetAllPostsMetaData() []Posts {
	var posts []Posts
	db.Raw(`SELECT post_id,post_title,author_name,post_likes FROM posts ORDER BY post_likes DESC`).Scan(&posts)
	return posts
}
func CheckPostTitleExists(posttitle string) (bool, error) {
	var exists bool
	r := db.Raw("SELECT EXISTS (SELECT 1 FROM posts WHERE post_title = ?)", posttitle).Scan(&exists)
	if r.Error != nil {
		return exists, r.Error
	}
	return exists, nil
}
func UpdatePostTitle(postid uint64, posttitle string) error {
	tx := db.Begin()
	r := tx.Raw("UPDATE posts SET post_title=? WHERE postid=?", posttitle, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
		return nil
	}
}
func UpdatePostContent(postid uint64, postContent string) error {
	tx := db.Begin()
	r := tx.Raw("UPDATE posts SET post_content=? WHERE postid=?", postContent, postid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
		return nil
	}

}
func GetPostsByAuthorID(authorid uint64, limit uint64, offset uint64) []Posts {
	var posts []Posts
	db.Raw(`SELECT 
	post_id,post_title,author_name,post_likes,createdat 
	FROM posts WHERE author_id=? 
	ORDER BY post_likes DESC LIMIT ? OFFSET ?`, authorid, limit, offset*limit).Scan(&posts)
	return posts
}

func GetPostsMetaData(offset uint64, limit uint64) []Posts {
	var posts []Posts
	db.Raw(`SELECT post_id,post_title,
	author_name,post_likes 
	FROM posts ORDER BY post_likes DESC LIMIT ? OFFSET ?`, limit, offset*limit).Scan(&posts)
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
