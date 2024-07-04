package models

import "time"

type Users struct {
	UserID    int64     `gorm:"primaryKey"`
	Username  string    `db:"username"`
	About     string    `db:"about"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	ImageURL  string    `db:"imageURL"`
	Active    bool      `db:"active"`
	Createdat time.Time `db:"createdAt"`
	Updatedat time.Time `db:"updatedAt"`
}

type Posts struct {
	Post_id       uint64    `db:"post_id"`
	Author_id     uint64    `db:"author_id"`
	Author_name   string    `db:"author_name"`
	Post_title    string    `db:"post_title"`
	Post_content  string    `db:"post_content"`
	Post_likes    uint32    `db:"post_likes"`
	Createdat     time.Time `db:"createdat"`
	Createdat_str string
	// UpdatedAt    time.Time `db:"updatedat"`
}

type Comment struct {
	Comment_id      uint64    `db:"comment_id"`
	Post_id         uint64    `db:"post_id"`
	User_id         uint64    `db:"user_id"`
	Username        string    `db:"username"`
	Comment_content string    `db:"comment_content"`
	Comment_likes   uint64    `db:"comment_likes"`
	Createdat       time.Time `db:"createdat"`
	Createdat_str   string
	Updatedat       time.Time `db:"updatedat"`
	Updatedat_str   string
}

type CommentsWithReactions struct {
	Comment_id      uint64    `db:"comment_id"`
	Post_id         uint64    `db:"post_id"`
	User_id         uint64    `db:"user_id"`
	Username        string    `db:"username"`
	Comment_content string    `db:"comment_content"`
	Comment_likes   uint64    `db:"comment_likes"`
	Createdat       time.Time `db:"createdat"`
	Createdat_str   string
	Updatedat       time.Time `db:"updatedat"`

	Liked    bool `db:"user_liked"`
	Disliked bool `db:"user_disliked"`
}

type UserAndPost struct {
	User  Users
	Posts []Posts
}
type UsernameAndPost struct {
	Username string
	Post     Posts
}

type PostComments_WithUserPreference struct {
	Post               Posts
	PostLikedByUser    bool
	PostDislikedByUser bool
	Comments           []CommentsWithReactions
}
type PostandComments struct {
	Post     Posts
	Comments []Comment
}
type UsernameAndComment struct {
	User_id         uint64
	Username        string
	Comment_id      uint64
	Comment_content string
	Comment_likes   uint64
	CreatedAt       string
}
type Passanduserid struct {
	Password string `db:"password"`
	User_id  uint64 `db:"user_id"`
}
type PostAndUserPreferences struct {
	Post               Posts
	Username           string
	PostLikedByUser    bool
	PostDislikedByUser bool
}

type UserInfo struct {
	UserID   int64
	Username string
	About    string
	JoinDate string
}
