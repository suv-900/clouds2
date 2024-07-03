package models

func LikeComment(commentid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("INSERT INTO comments_liked_by_users(user_id,comment_id) VALUES(?,?) ", userid, commentid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE comments SET comment_likes = comment_likes + 1 WHERE comment_id=?", commentid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}

func DislikeComment(commentid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("INSERT INTO comments_disliked_by_users(user_id,comment_id) VALUES(?,?) ", userid, commentid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE comments SET comment_likes = comment_likes - 1 WHERE comment_id=?", commentid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}
func RemoveLikeFromComment(commentid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("DELETE FROM comments_liked_by_users WHERE user_id=? AND comment_id=?", userid, commentid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE comments SET comment_likes=comment_likes - 1 WHERE comment_id=?", commentid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}

func RemoveDislikeFromComment(commentid uint64, userid uint64) error {
	tx := db.Begin()
	r := tx.Exec("DELETE FROM comments_disliked_by_users WHERE user_id=? AND comment_id=?", userid, commentid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}

	tx = db.Begin()
	r = tx.Exec("UPDATE comments SET comment_likes=comment_likes + 1 WHERE comment_id=?", commentid)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	} else {
		tx.Commit()
	}
	return nil
}

func GetUserCommentReaction(postid uint64, userid uint64) []CommentsWithReactions {
	var comments []CommentsWithReactions
	sql := `SELECT c.*, 
		CASE 
			WHEN l.like_id IS NOT NULL THEN true
			ELSE false
		END AS user_liked,
		CASE 
			WHEN d.dislike_id IS NOT NULL THEN true
			ELSE false
		END AS user_disliked
		FROM
			comments c
		LEFT JOIN
			comments_liked_by_users l
		ON 
			c.comment_id = l.comment_id AND l.user_id = ?
		LEFT JOIN
			comments_disliked_by_users d 
		ON 
			c.comment_id = d.comment_id AND d.user_id = ?
		WHERE
			c.post_id = ?
		ORDER BY
			c.comment_likes
		DESC LIMIT 5
		`
	db.Raw(sql, userid, userid, postid).Scan(&comments)
	return comments
}

func GetAllCommentsByPostID(postid uint64) []Comment {
	var comments []Comment
	sql := "SELECT * FROM comments WHERE post_id=? ORDER BY comment_likes DESC "
	db.Raw(sql, postid).Scan(&comments)
	return comments
}
func GetComments(postid uint64, limit uint64, offset uint64) ([]Comment, error) {
	var comments []Comment
	sql := "SELECT * FROM comments WHERE post_id=? ORDER BY comment_likes DESC LIMIT ? OFFSET ?"
	r := db.Raw(sql, postid, limit, offset*limit).Scan(&comments)
	if r.Error != nil {
		return comments, r.Error
	}
	return comments, nil
}
func AddComment(postid uint64, userid uint64, username string, comment_content string) (uint64, error) {
	var commentID uint64

	tx := db.Begin()
	sql := "INSERT INTO comments (user_id,post_id,username,comment_content,comment_likes) VALUES(?,?,?,?,?) RETURNING comment_id"

	r := tx.Raw(sql, userid, postid, username, comment_content, 0).Scan(&commentID)

	if r.Error != nil {
		tx.Rollback()
		return commentID, r.Error
	} else {
		tx.Commit()
	}
	return commentID, nil
}

func EditComment(commentId uint64, comment string) {
	db.Exec("UPDATE comments SET comment_content=? WHERE comment_id=?", commentId, comment)
}

func RemoveComment(commentId uint64) {
	db.Exec("DELETE * FROM comments WHERE comment_id=?", commentId)
}
