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

	router.GET("/subscribe/:email:uuid", func(c *gin.Context) {
		email := c.Param("email")
		uuid := c.Param("uuid")
		_, err := services.ActivateSubscription(uuid, email)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Status": email + " has been activated",
		})
	})

	router.GET("/subscribe/?email=:email?uuid=:uuid", func(c *gin.Context) {
		email := c.Param("email")
		_, err := services.Subscribe(email)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Status": "Awesome. We added " + email + " to our database. W8 for best newsletter ever!",
		})
	})
}
