package main

import (
	"MailNews.Subscriber/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	ginHost := gin.Default()
	routes.SubscriberRoute(ginHost)
	err := ginHost.Run()
	if err != nil {
		return
	}
	fmt.Println("Test")
}
