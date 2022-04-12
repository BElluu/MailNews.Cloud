package routes

import (
	"MailNews.Subscriber/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SubscriberRoute(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/subscribe/:email/:uuid", func(c *gin.Context) {
		email := c.Param("email")
		uuid := c.Param("uuid")
		_, err := services.ActivateSubscription(uuid, email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Status": email + " has been activated",
		})
	})

	router.GET("/subscribe/:email", func(c *gin.Context) {
		email := c.Param("email")
		_, err := services.Subscribe(email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Status": "Awesome. We added " + email + " to our database. W8 for best newsletter ever!",
		})
	})
}
