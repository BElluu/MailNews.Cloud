package main

import (
	"MailNews.Subscriber/database"
	"MailNews.Subscriber/routes"
	"MailNews.Subscriber/services"
	"github.com/gin-gonic/gin"
)

func main() {

	client := database.CreateLocalClient()
	database.PrepareDatabaseTables(client)
	services.FetchFeeds()

	ginHost := gin.Default()
	routes.SubscriberRoute(ginHost)
	err := ginHost.Run()
	if err != nil {
		return
	}
}
