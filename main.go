package main

import (
	"MailNews.Subscriber/common"
	"MailNews.Subscriber/database"
	"MailNews.Subscriber/routes"
	"MailNews.Subscriber/services"
	"github.com/gin-gonic/gin"
)

func main() {

	client := common.CreateLocalClient()
	database.PrepareDatabaseTables(client)
	database.PrintAllTables(client)
	services.FetchFeeds()
	ginHost := gin.Default()
	routes.SubscriberRoute(ginHost)
	err := ginHost.Run()
	if err != nil {
		return
	}
}
