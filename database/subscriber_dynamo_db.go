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

func ListTables(d *dynamodb.Client) {
	tables, err := d.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatal("ListTables failed", err)
	}
	for wtf := range tables.TableNames {
		println(wtf)
	}
}

func tableExists(d *dynamodb.Client, name string) bool {
	tables, err := d.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
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

func CreateTableIfNotExists(d *dynamodb.Client, tableName string) {
	if tableExists(d, tableName) {
		log.Printf("table=%v already exists\n", tableName)
		return
	}
	_, err := d.CreateTable(context.TODO(), buildCreateTableInput(tableName))
	if err != nil {
		log.Fatal("CreateTable failed", err)
	}
	log.Printf("created table=%v\n", tableName)
}

func buildCreateTableInput(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
	}
}
