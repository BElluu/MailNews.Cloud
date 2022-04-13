package main

import (
	"MailNews.Subscriber/database"
	"MailNews.Subscriber/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	client := database.CreateLocalClient()
	database.CreateTableIfNotExists(client, "Suby")

	ginHost := gin.Default()
	routes.SubscriberRoute(ginHost)
	err := ginHost.Run()
	if err != nil {
		return
	}

	// Create local table
	/*	fmt.Println("Test")
		client := database.CreateLocalClient()
		database.CreateTableIfNotExists(client, "Suby")

		// List tables
		database.ListTables(client)
		fmt.Println("Koniec")

		//add item
		var sub = models.Subscriber{
			SubscriberId:   2,
			Email:          "paulajakubasik@gmail.com",
			ActivateURL:    "http://www.devopsowy.pl/actiave/?email=pjakubasik",
			UnSubscribeURL: "http://www.devopsowy.pl/unsubscribe/?email=pjakubasik",
			IsActive:       true,
		}
		database.AddItem(context.TODO(), sub, client, "Suby")*/
}

// aws dynamodb list-tables --endpoint-url http://127.0.0.1:8000 --region local
