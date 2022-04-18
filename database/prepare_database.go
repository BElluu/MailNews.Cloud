package database

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"time"
)

func CreateLocalClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("local"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://127.0.0.1:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "HardTesty",
			},
		}),
	)
	if err != nil {
		panic(err)
	}
	return dynamodb.NewFromConfig(cfg)
}

func PrepareDatabaseTables(client *dynamodb.Client) {
	createTableSubscribersIfNotExists(client)
	createTableFeedsIfNotExists(client)
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

	fillConfigTable(context.Background(), client)
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
func fillConfigTable(ctx context.Context, client *dynamodb.Client) {
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

	_, err := svc.PutItem(ctx, input)
	if err != nil {
		log.Fatalf("%s", err)
	}
	input2 := &dynamodb.PutItemInput{
		Item:      configMap2,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(ctx, input2)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
