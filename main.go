package main

import (
	"fmt"
	"log"
	"ticketing/cmd"
	"ticketing/common"
	"ticketing/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	docs.SwaggerInfo.Title = common.Configuration.Title
	docs.SwaggerInfo.Schemes = common.Configuration.Schemes
	docs.SwaggerInfo.Host = common.Configuration.Host
	docs.SwaggerInfo.BasePath = common.Configuration.BasePath

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	ticket := router.Group("/ticket")
	ticket.POST("", cmd.CreateTicket)
	ticket.GET(":ticketId", cmd.GetTicket)

	comment := router.Group("/comment")
	comment.GET(":key", cmd.GetComments)
	comment.POST("", cmd.CreateComment)
	comment.POST(":commentId/react", cmd.ReactComment)

	admin := router.Group("/admin")

	adminComment := admin.Group("comment")
	adminComment.POST(":commentId", cmd.ConfirmComment)
	adminComment.POST(":commentId/reject", cmd.RejectComment)

	adminTicket := admin.Group("ticket")
	adminTicket.POST(":ticketId", cmd.CloseTicket)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(fmt.Sprintf("localhost:%d", common.Configuration.ListenPort))
	if err != nil {
		log.Fatal(err)
	}
}
