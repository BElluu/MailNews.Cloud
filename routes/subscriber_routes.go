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

	router.GET("/activate", func(c *gin.Context) {
		email := c.Query("email")
		uuid := c.Query("uuid")
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

	router.GET("/subscribe", func(c *gin.Context) {
		email := c.Query("email")
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

	router.GET("/unsubscribe", func(c *gin.Context) {
		email := c.Query("email")
		uuid := c.Query("uuid")
		_, err := services.UnSubscribe(uuid, email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Status": email + " unsubscribed",
		})
	})
}
