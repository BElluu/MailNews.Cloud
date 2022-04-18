package main

import (
	"MailNews.Subscriber/routes"
	"MailNews.Subscriber/rss"
	"github.com/gin-gonic/gin"
)

func main() {

	/*	client := database.CreateLocalClient()
		database.CreateTableSubscribersIfNotExists(client)*/
	rss.AmazonRSSParser()

	ginHost := gin.Default()
	routes.SubscriberRoute(ginHost)
	err := ginHost.Run()
	if err != nil {
		return
	}
}
