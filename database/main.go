package main

import (
	"MailNews.Cloud/backend/common"
	dbservice "MailNews.Cloud/database/services"
)

func main() {
	client := common.CreateLocalClient()
	dbservice.PrepareDatabaseTables(client)
	dbservice.PrintAllTables(client)
}
