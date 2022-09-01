package app

import (
	"fmt"
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
	IsConfirmed     *bool
	LikeCount       int `gorm:"<-"`
	DislikeCount    int `gorm:"<-"`
	ConfirmUserId   *int
	ConfirmDateTime *time.Time
	ReplyToId       *uint
	ReplyTo         *Comment `gorm:"foreignkey:ReplyToId"`
}

type React struct {
	ID        uint
	Type      int
	UserId    int
	ClientId  string
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

func (comment *Comment) changeConfirmed(userId int, isConfirmed bool) {
	now := time.Now()
	comment.ConfirmDateTime = &now
	comment.IsConfirmed = &isConfirmed
	comment.ConfirmUserId = &userId
}

func (comment *Comment) ToModel() CommentModel {
	return CommentModel{
		Id:           comment.ID,
		UserId:       comment.UserId,
		UserName:     comment.UserName,
		Content:      comment.Content,
		CreatedAt:    comment.CreatedAt,
		LikeCount:    comment.LikeCount,
		DisLikeCount: comment.DislikeCount,
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

	comment.changeConfirmed(userId, true)

	db.Updates(&comment)
	return comment.ToModel(), nil
}

func RejectComment(commentId uint, userId int, db *gorm.DB) (CommentModel, error) {
	comment := Comment{}
	db.First(&comment, commentId)
	if db.Error != nil {
		return comment.ToModel(), db.Error
	}

	comment.changeConfirmed(userId, false)

	db.Updates(&comment)
	print(db.Error)
	return comment.ToModel(), nil
}

func GetComment(Id uint, db *gorm.DB) CommentModel {
	comment := Comment{}

	subQuery := db.Table("React").
		Select("CommentId, "+
			"SUM(CASE WHEN Type=? THEN 1 ELSE 0 END) AS LikeCount,"+
			"SUM(CASE WHEN Type=? THEN 1 ELSE 0 END) AS DislikeCount", Like, Dislike).
		Group("CommentId")

	db.Model(&Comment{}).
		Select("*").
		Joins("LEFT JOIN (?) T ON T.CommentId = Comment.ID", subQuery).
		Where("Comment.ID = ?", Id).
		Find(&comment)

	return comment.ToModel()
}

func GetComments(key string, db *gorm.DB) []CommentModel {
	comments := make([]Comment, 0)
	isConfirmed := true

	subQuery := db.Table("React").
		Select("CommentId, "+
			"SUM(CASE WHEN Type=? THEN 1 ELSE 0 END) AS LikeCount,"+
			"SUM(CASE WHEN Type=? THEN 1 ELSE 0 END) AS DislikeCount", Like, Dislike).
		Group("CommentId")

	db.Model(&Comment{}).
		Select("*").
		Joins("LEFT JOIN (?) T ON T.CommentId = Comment.ID", subQuery).
		Where(Comment{Key: key, IsConfirmed: &isConfirmed}).
		Find(&comments)

	return convertComments(&comments)
}

func CreateReact(commentId uint, userId int, clientId string, reactType int, db *gorm.DB) (CommentModel, error) {
	comment := Comment{}
	db.Find(&comment, commentId)

	if comment.ID == 0 {
		return comment.ToModel(), fmt.Errorf("کامنت پیدا نشد")
	}

	oldReact := React{}
	db.Model(&oldReact).Where(React{CommentId: commentId, UserId: userId}).
		Or(React{CommentId: commentId, ClientId: clientId}).
		Find(&oldReact)

	if oldReact.ID != 0 && reactType == oldReact.Type {
		return comment.ToModel(), fmt.Errorf("شما قبلا نظر داده اید")
	} else if oldReact.ID != 0 {
		oldReact.Type = reactType
		db.Updates(&oldReact)

		return GetComment(commentId, db), nil
	}

	db.Create(&React{
		Type:      reactType,
		UserId:    userId,
		ClientId:  clientId,
		CommentId: commentId,
		Comment:   comment,
	})

	return GetComment(commentId, db), nil
}
