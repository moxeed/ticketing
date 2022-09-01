package app

import (
	"gorm.io/gorm"
	"time"
)

const (
	Like = iota + 1
	Dislike
)

type Comment struct {
	gorm.Model
	Key             string
	UserId          *int
	UserName        string
	Content         string
	IsConfirmed     bool
	ConfirmUserId   *int
	ConfirmDateTime *time.Time
	ReplyToId       *uint
	ReplyTo         *Comment `gorm:"foreignkey:ReplyToId"`
}

type React struct {
	ID        int
	Type      int
	CommentId uint
	Comment   Comment `gorm:"foreignkey:CommentId"`
}

func NewComment(key string, userId *int, userName string, content string, replyToId *uint) Comment {
	return Comment{
		Key:       key,
		UserId:    userId,
		UserName:  userName,
		Content:   content,
		ReplyToId: replyToId,
	}
}

func (comment *Comment) Confirm(userId int) {
	now := time.Now()
	comment.ConfirmDateTime = &now
	comment.IsConfirmed = true
	comment.ConfirmUserId = &userId
}

func (comment *Comment) ToModel() CommentModel {
	return CommentModel{
		ID:           comment.ID,
		UserId:       comment.UserId,
		UserName:     comment.UserName,
		Content:      comment.Content,
		CreatedAt:    comment.CreatedAt,
		LikeCount:    0,
		DisLikeCount: 0,
		ReplyToId:    comment.ReplyToId,
	}
}

func convertComments(tickets *[]Comment) []CommentModel {
	models := make([]CommentModel, 0)

	for _, ticket := range *tickets {
		models = append(models, ticket.ToModel())
	}

	return models
}

func CreateComment(model CommentCreateModel, db *gorm.DB) CommentModel {
	comment := NewComment(model.Key,
		model.UserId,
		model.UserName,
		model.Content,
		model.ReplyToId)

	db.Create(&comment)

	return comment.ToModel()
}

func ConfirmComment(commentId uint, userId int, db *gorm.DB) (CommentModel, error) {
	comment := Comment{}
	db.First(&comment, commentId)
	if db.Error != nil {
		return comment.ToModel(), db.Error
	}

	comment.Confirm(userId)
	return comment.ToModel(), nil
}

func GetComments(key string, db *gorm.DB) []CommentModel {
	comments := make([]Comment, 0)
	db.Where(Comment{Key: key, IsConfirmed: true}).Find(&comments)

	return convertComments(&comments)
}
