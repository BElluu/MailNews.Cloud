package database

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
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

func createTableSubscribersIfNotExists(d *dynamodb.Client) {
	if tableExists(d, SubscriberTable) {
		log.Printf("table=%v already exists\n", SubscriberTable)
		return
	}
	_, err := d.CreateTable(context.Background(), buildCreateTableInputSubscribers())
	if err != nil {
		log.Fatal("CreateTable failed", err)
	}
	log.Printf("created table=%v\n", SubscriberTable)
}

func createTableFeedsIfNotExists(d *dynamodb.Client) {
	if tableExists(d, FeedTable) {
		log.Printf("table=%v already exists\n", FeedTable)
		return
	}
	_, err := d.CreateTable(context.Background(), buildCreateTableInputFeeds())
	if err != nil {
		log.Fatal("CreateTable failed", err)
	}
	log.Printf("created table=%v\n", FeedTable)
}

func createTableConfigIfNotExists(d *dynamodb.Client) {
	if tableExists(d, "MailNewsConfig") {
		log.Printf("table=%v already exists\n", "MailNewsConfig")
		return
	}
	_, err := d.CreateTable(context.Background(), buildCreateTableInputConfiguration())
	if err != nil {
		log.Fatal("CreateTable failed", err)
	}
	log.Printf("created table=%v\n", "MailNewsConfig")
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
