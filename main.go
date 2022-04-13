package main

import (
	"MailNews.Subscriber/database"
	"MailNews.Subscriber/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	client := database.CreateLocalClient()
	database.CreateTableSubscribersIfNotExists(client)

	ginHost := gin.Default()
	routes.SubscriberRoute(ginHost)
	err := ginHost.Run()
	if err != nil {
		return
	}
}
