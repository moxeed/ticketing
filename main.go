package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"ticketing/cmd"
	"ticketing/docs"
)

func main() {

	docs.SwaggerInfo.Title = "ticketing"
	router := gin.Default()

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

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}
