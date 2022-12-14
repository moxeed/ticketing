package cmd

import (
	"strconv"
	"ticketing/app"
	"ticketing/common"

	"github.com/gin-gonic/gin"
)

// CreateComment godoc
// @Tags         comment
// @Description Create New Comment
// @Accept      json
// @Produce     json
// @Param       ticket body     app.CommentCreateModel true "comment data"
// @Success     200    {object} app.CommentModel
// @Failure     400    {object} common.Error
// @Failure     404    {object} common.Error
// @Router      /comment [post]
func CreateComment(ctx *gin.Context) {
	model := app.CommentCreateModel{}
	if err := ctx.BindJSON(&model); err != nil {
		ctx.JSON(400, common.Error{Error: err.Error(), Status: 400})
		return
	}

	result := app.CreateComment(model, common.Db)
	ctx.JSON(200, result)
}

// ConfirmComment godoc
// @Tags         comment
// @Description Confirm Ticket
// @Accept      json
// @Produce     json
// @Param       commentId path     uint true "comment id"
// @Param       userId    query    int  true "user id"
// @Success     200       {object} app.CommentModel
// @Failure     400       {object} common.Error
// @Failure     404       {object} common.Error
// @Router      /admin/comment/{commentId} [post]
func ConfirmComment(ctx *gin.Context) {
	commentId, err := strconv.ParseUint(ctx.Param("commentId"), 10, 0)
	if err != nil {
		ctx.JSON(400, common.Error{Error: "شناسه کامنت نادرست است", Status: 400})
		return
	}

	userId, err := strconv.ParseInt(ctx.Query("userId"), 10, 0)
	if err != nil {
		ctx.JSON(400, common.Error{Error: "شناسه کاربر نادرست است", Status: 400})
		return
	}

	result, err := app.ConfirmComment(uint(commentId), int(userId), common.Db)

	if err != nil {
		ctx.JSON(400, common.Error{Error: err.Error(), Status: 400})
		return
	}

	ctx.JSON(200, result)
}

// RejectComment godoc
// @Tags         comment
// @Description Confirm Ticket
// @Accept      json
// @Produce     json
// @Param       commentId path     uint true "comment id"
// @Param       userId    query    int  true "user id"
// @Success     200       {object} app.CommentModel
// @Failure     400       {object} common.Error
// @Failure     404       {object} common.Error
// @Router      /admin/comment/{commentId}/reject [post]
func RejectComment(ctx *gin.Context) {
	commentId, err := strconv.ParseUint(ctx.Param("commentId"), 10, 0)
	if err != nil {
		ctx.JSON(400, common.Error{Error: "شناسه کامنت نادرست است", Status: 400})
		return
	}

	userId, err := strconv.ParseInt(ctx.Query("userId"), 10, 0)
	if err != nil {
		ctx.JSON(400, common.Error{Error: "شناسه کاربر نادرست است", Status: 400})
		return
	}

	result, err := app.RejectComment(uint(commentId), int(userId), common.Db)

	if err != nil {
		ctx.JSON(400, common.Error{Error: err.Error(), Status: 400})
		return
	}

	ctx.JSON(200, result)
}

// ReactComment godoc
// @Tags        comment
// @Description Confirm Ticket
// @Accept      json
// @Produce     json
// @Param       commentId path     uint    true "comment id"
// @Param       userId    query    int     true "user id"
// @Param       clientId  query    string  true "client id"
// @Param       reactType query    int     true "react type"
// @Success     200       {object} app.CommentModel
// @Failure     400       {object} common.Error
// @Failure     404       {object} common.Error
// @Router      /comment/{commentId}/react [post]
func ReactComment(ctx *gin.Context) {
	commentId, err := strconv.ParseUint(ctx.Param("commentId"), 10, 0)
	if err != nil {
		ctx.JSON(400, common.Error{Error: "شناسه کامنت نادرست است", Status: 400})
		return
	}

	userId, err := strconv.ParseInt(ctx.Query("userId"), 10, 0)
	if err != nil {
		ctx.JSON(400, common.Error{Error: "شناسه کاربر نادرست است", Status: 400})
		return
	}

	reactType, err := strconv.ParseInt(ctx.Query("reactType"), 10, 0)
	if err != nil {
		ctx.JSON(400, common.Error{Error: "شناسه کاربر نادرست است", Status: 400})
		return
	}

	clientId := ctx.Query("clientId")

	result, err := app.CreateReact(uint(commentId), int(userId), clientId, int(reactType), common.Db)

	if err != nil {
		ctx.JSON(400, common.Error{Error: err.Error(), Status: 400})
		return
	}

	ctx.JSON(200, result)
}

// GetComments godoc
// @Tags         comment
// @Description Confirm Ticket
// @Accept      json
// @Produce     json
// @Param       key path     string true "comment group key"
// @Success     200 {array}  app.CommentModel
// @Failure     400 {object} common.Error
// @Failure     404 {object} common.Error
// @Router      /comment/{key} [get]
func GetComments(ctx *gin.Context) {
	key := ctx.Param("key")

	result := app.GetComments(key, common.Db)

	ctx.JSON(200, result)

}

// GetUserComments godoc
// @Tags        comment
// @Description Get User Ticket
// @Accept      json
// @Produce     json
// @Param       userId path  int true "User Id"
// @Success     200 {array}  app.CommentModel
// @Failure     400 {object} common.Error
// @Failure     404 {object} common.Error
// @Router      /comment/user/{userId} [get]
func GetUserComments(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 0)
	if err != nil {
		ctx.JSON(400, common.Error{Error: "شناسه کاربر نادرست است", Status: 400})
		return
	}

	result := app.GetUserComments(int(userId), common.Db)

	ctx.JSON(200, result)
}
