package app

import "time"

type CommentCreateModel struct {
	Key       string
	UserId    *int
	UserName  string
	Content   string
	ReplyToId *uint
}

type CommentModel struct {
	Id           uint
	UserId       *int
	UserName     string
	Content      string
	CreatedAt    time.Time
	LikeCount    int
	DisLikeCount int
	ReplyToId    *uint
}
