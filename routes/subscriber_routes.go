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

	router.GET("/subscribe/:email", func(c *gin.Context) {
		email := c.Param("email")
		_, err := services.Subscribe(email)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			c.JSON(200, gin.H{
				"Status": "Awesome. W8 for best newsletter ever!",
			})
			return
		}
	})
}
