package main

import (
	"MailNews.Subscriber/database"
	"fmt"
)

func main() {
	/*	ginHost := gin.Default()
		routes.SubscriberRoute(ginHost)
		err := ginHost.Run()
		if err != nil {
			return
		}*/
	fmt.Println("Test")
	client := database.CreateLocalClient()
	database.CreateTableIfNotExists(client, "wtasf")
	database.ListTables(client)
	fmt.Println("Koniec")
}
// aws dynamodb list-tables --endpoint-url http://127.0.0.1:8000 --region local