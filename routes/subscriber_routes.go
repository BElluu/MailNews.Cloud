package routes

import (
	"MailNews.Subscriber/services"
	"github.com/gin-gonic/gin"
)

func SubscriberRoute(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	email := "bkomendarczuk@gmail.com"
	router.GET("/test", func(c *gin.Context) {
		_, err := services.Subscribe(email)
		if err != nil {
			return
		}
	})
}
