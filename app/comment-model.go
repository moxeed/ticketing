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
	Id           uint      `json:"id,omitempty"`
	UserId       *int      `json:"userId,omitempty"`
	UserName     string    `json:"userName,omitempty"`
	Content      string    `json:"content,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	LikeCount    int       `json:"likeCount,omitempty"`
	DisLikeCount int       `json:"disLikeCount,omitempty"`
	ReplyToId    *uint     `json:"replyToId,omitempty"`
}
