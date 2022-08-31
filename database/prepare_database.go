package database

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"time"
)

func PrepareDatabaseTables(client *dynamodb.Client) {
	createTableSubscribersIfNotExists(client)
	//createTableFeedsIfNotExists(client) // split to providers
	createProvidersTablesIfNotExists(client)
	createTableConfigIfNotExists(client)
}

func createTableSubscribersIfNotExists(client *dynamodb.Client) {
	if tableExists(client, SubscriberTable) {
		log.Printf("table=%v already exists\n", SubscriberTable)
		return
	}
	_, err := client.CreateTable(context.Background(), buildCreateTableInputSubscribers())
	if err != nil {
		log.Fatal("CreateTable failed", err)
	}
	log.Printf("created table=%v\n", SubscriberTable)
}

func createTableFeedsIfNotExists(client *dynamodb.Client) {
	if tableExists(client, FeedTable) {
		log.Printf("table=%v already exists\n", FeedTable)
		return
	}
	_, err := client.CreateTable(context.Background(), buildCreateTableInputFeeds())
	if err != nil {
		log.Fatal("CreateTable failed", err)
	}
	log.Printf("created table=%v\n", FeedTable)
}

func createProvidersTablesIfNotExists(client *dynamodb.Client) {
	tablesToCreate := prepareProviderTablesDefinition()

	for tableName, table := range tablesToCreate {
		if tableExists(client, tableName) {
			log.Printf("table=%v already exists\n", tableName)
			return
		}

		_, err := client.CreateTable(context.Background(), &table)
		if err != nil {
			log.Fatal("CreateTable failed", err)
		}
		log.Printf("created table=%v\n", tableName)
	}
}

func createTableConfigIfNotExists(client *dynamodb.Client) {
	if tableExists(client, "MailNewsConfig") {
		log.Printf("table=%v already exists\n", "MailNewsConfig")
		return
	}
	_, err := client.CreateTable(context.Background(), buildCreateTableInputConfiguration())
	if err != nil {
		log.Fatal("CreateTable failed", err)
	}
	log.Printf("created table=%v\n", "MailNewsConfig")

	fillConfigTable(client)
}

func ListTables(d *dynamodb.Client) {
	tables, err := d.ListTables(context.Background(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatal("ListTables failed", err)
	}
	for wtf := range tables.TableNames {
		println(wtf)
	}
}

func buildCreateTableInputSubscribers() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("UUID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("Email"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("UUID"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("Email"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(SubscriberTable),
		BillingMode: types.BillingModePayPerRequest,
	}
}

func buildCreateTableInputFeeds() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("UUID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("UUID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(FeedTable),
		BillingMode: types.BillingModePayPerRequest,
	}
}

func prepareProviderTablesDefinition() map[string]dynamodb.CreateTableInput {
	providers := []string{AWSTable, AzureTable, GCPTable}
	tables := make(map[string]dynamodb.CreateTableInput)

	for _, provider := range providers {
		tables[provider] =
			dynamodb.CreateTableInput{
				AttributeDefinitions: []types.AttributeDefinition{
					{
						AttributeName: aws.String("UUID"),
						AttributeType: types.ScalarAttributeTypeS,
					},
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("UUID"),
						KeyType:       types.KeyTypeHash,
					},
				},
				TableName:   aws.String(provider),
				BillingMode: types.BillingModePayPerRequest,
			}
	}
	return tables
}

func buildCreateTableInputConfiguration() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Name"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String("MailNewsConfig"),
		BillingMode: types.BillingModePayPerRequest,
	}
}

func tableExists(d *dynamodb.Client, name string) bool {
	tables, err := d.ListTables(context.Background(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatal("ListTables failed", err)
	}
	for _, n := range tables.TableNames {
		if n == name {
			return true
		}
	}
	return false
}
func fillConfigTable(client *dynamodb.Client) {
	svc := client
	tableName := "MailNewsConfig"

	configMap := map[string]types.AttributeValue{
		"Name":  &types.AttributeValueMemberS{Value: "LastFetchFeedsDate"},
		"Value": &types.AttributeValueMemberS{Value: time.Now().Format("02-01-2006 15:01:05")},
	}
	configMap2 := map[string]types.AttributeValue{
		"Name":  &types.AttributeValueMemberS{Value: "LastSendMailDate"},
		"Value": &types.AttributeValueMemberS{Value: time.Now().Format("02-01-2006 15:01:05")},
	}

	input := &dynamodb.PutItemInput{
		Item:      configMap,
		TableName: aws.String(tableName),
	}

	_, err := svc.PutItem(context.Background(), input)
	if err != nil {
		log.Fatalf("%s", err)
	}
	input2 := &dynamodb.PutItemInput{
		Item:      configMap2,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(context.Background(), input2)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
