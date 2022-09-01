package cmd

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"ticketing/app"
	"ticketing/common"
)

// CreateComment godoc
// @Description  Create New Comment
// @Accept       json
// @Produce      json
// @Param        ticket body app.CommentCreateModel true "comment data"
// @Success      200  {object}  app.CommentModel
// @Failure      400  {object}  common.Error
// @Failure      404  {object}  common.Error
// @Router       /comment [post]
func CreateComment(ctx *gin.Context) {
	model := app.CommentCreateModel{}
	if err := ctx.BindJSON(&model); err != nil {
		ctx.JSON(400, err)
		return
	}

	db := common.OpenDb()

	result := app.CreateComment(model, db)
	ctx.JSON(200, result)
}

// ConfirmComment godoc
// @Description  Confirm Ticket
// @Accept       json
// @Produce      json
// @Param        commentId path uint true "comment id"
// @Param        userId query int true "user id"
// @Success      200  {object}  app.CommentModel
// @Failure      400  {object}  common.Error
// @Failure      404  {object}  common.Error
// @Router       /comment/{commentId} [post]
func ConfirmComment(ctx *gin.Context) {
	commentId, err := strconv.ParseUint(ctx.Param("commentId"), 10, 0)
	if err != nil {
		ctx.JSON(400, common.Error{Error: "شناسه کامنت نادرست است", Status: 400})
	}

	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 0)
	if err != nil {
		ctx.JSON(400, common.Error{Error: "شناسه کاربر نادرست است", Status: 400})
	}

	db := common.OpenDb()
	result, err := app.ConfirmComment(uint(commentId), int(userId), db)

	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(200, result)
}

// GetComments godoc
// @Description  Confirm Ticket
// @Accept       json
// @Produce      json
// @Param        key path string true "comment group key"
// @Success      200  {array}  app.CommentModel
// @Failure      400  {object}  common.Error
// @Failure      404  {object}  common.Error
// @Router       /comment/{key} [get]
func GetComments(ctx *gin.Context) {
	key := ctx.Param("key")

	db := common.OpenDb()
	result := app.GetComments(key, db)

	ctx.JSON(200, result)
}
