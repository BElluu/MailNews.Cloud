package main

import (
	"MailNews.Cloud/backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	ginHost := gin.Default()
	routes.SubscriberRoute(ginHost)
	err := ginHost.Run()
	if err != nil {
		return
	}
}
