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

	router.POST("/ticket", cmd.CreateTicket)
	router.POST("/ticket/:ticketId", cmd.CloseTicket)

	router.POST("/tickets", cmd.GetTickets)
	router.GET("/ticket/:ticketId", cmd.GetTicket)

	router.GET("/comment/:key", cmd.GetComments)
	router.POST("/comment", cmd.CreateComment)
	router.POST("/admin/comment/:commentId", cmd.ConfirmComment)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}
