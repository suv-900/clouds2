package models

// call responsively
func LikeAComment(commentid uint64) {
	db.Exec("UPDATE comments SET comment_likes=comment_likes+1 WHERE comment_id=?", commentid)
}

func DislikeAComment(commentid uint64) {
	db.Exec("UPDATE comments SET comment_likes=comment_likes-1 WHERE comment_id=?", commentid)
}

func Get5CommentsByPostID(postid uint64) ([]UsernameAndComment, error) {
	//commentarr := make([]UsernameAndComment, 5)
	commentsvec := []UsernameAndComment{}
	//TODO OFFSET to hold a bar for next comments
	a := make(chan int, 1)
	var err error
	go func() {
		sql := `SELECT 
			comment_id,
			user_id,
			username,
			comment_content,
			comment_likes
			FROM comments WHERE post_id=? ORDER BY comment_likes DESC LIMIT 5 `
		r := db.Raw(sql, postid).Scan(&commentsvec)
		err = r.Error
		a <- 1
	}()
	<-a
	return commentsvec, err
	/*
	   		rawComment := UsernameAndComment{
	   			UserID:          comment.User_id,
	   			Username:        username,
	         CommentID:       comment.Comment_id,
	   			Comment_content: comment.Comment_content,
	   		}
	   		commentarr = append(commentarr, rawComment)
	   	}
	*/
}

func GetAllCommentsByPostID(postid uint64) []UsernameAndComment {
	commentsvec := []UsernameAndComment{}
	//TODO OFFSET to hold a bar for next comments
	sql := "SELECT comment_id,user_id,username,comment_content,comment_likes FROM comments WHERE post_id=? ORDER BY comment_likes DESC "
	db.Raw(sql, postid).Scan(&commentsvec)
	return commentsvec
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
