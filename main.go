package main

import (
	"MailNews.Subscriber/database"
	"MailNews.Subscriber/models"
	"context"
	"fmt"
)

func main() {
	/*	ginHost := gin.Default()
		routes.SubscriberRoute(ginHost)
		err := ginHost.Run()
		if err != nil {
			return
		}*/

	// Create local table
	fmt.Println("Test")
	client := database.CreateLocalClient()
	database.CreateTableIfNotExists(client, "Suby")

	// List tables
	database.ListTables(client)
	fmt.Println("Koniec")

	//add item
	var sub = models.Subscriber{
		//SubId:          666,
		Email: "bkomendarczuk@gmail.com",
		/*		ActivateURL:    "http://www.devopsowy.pl/actiave/",
				UnSubscribeURL: "http://www.devopsowy.pl/unsubscribe/",*/
		/*		IsActive: false,*/
	}
	database.AddItem(context.TODO(), sub, client, "Suby")
}

// aws dynamodb list-tables --endpoint-url http://127.0.0.1:8000 --region local
