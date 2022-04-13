package database

import (
	"MailNews.Subscriber/models"
	"context"
	"errors"
	"fmt"
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
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
	}
}

func AddSubscriber(ctx context.Context, subscriber models.Subscriber, client *dynamodb.Client, table string) {
	svc := client
	tableName := table

	subscriberMap := map[string]types.AttributeValue{
		"UUID":           &types.AttributeValueMemberS{Value: subscriber.UUID},
		"Email":          &types.AttributeValueMemberS{Value: subscriber.Email},
		"ActivateURL":    &types.AttributeValueMemberS{Value: subscriber.ActivateURL},
		"UnSubscribeURL": &types.AttributeValueMemberS{Value: subscriber.UnSubscribeURL},
		"IsActive":       &types.AttributeValueMemberBOOL{Value: subscriber.IsActive},
	}

	input := &dynamodb.PutItemInput{
		Item:      subscriberMap,
		TableName: aws.String(tableName),
	}

	_, err := svc.PutItem(ctx, input)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Println("Added " + subscriber.Email + " to table" + table)
}

func DeleteSubscriber(ctx context.Context, uuid string, email string, client *dynamodb.Client, table string) (bool, error) {
	svc := client
	tableName := table

	subscriber := map[string]types.AttributeValue{
		"UUID":  &types.AttributeValueMemberS{Value: uuid},
		"Email": &types.AttributeValueMemberS{Value: email},
	}

	input := &dynamodb.DeleteItemInput{
		Key:       subscriber,
		TableName: aws.String(tableName),
	}

	_, err := svc.DeleteItem(ctx, input)
	if err != nil {
		return false, errors.New(err.Error())
	}
	return true, nil
}

func ActiveSubscriber(ctx context.Context, uuid string, email string, client *dynamodb.Client, table string) (bool, error) {
	svc := client
	tableName := table

	//
	key := map[string]types.AttributeValue{
		"UUID":  &types.AttributeValueMemberS{Value: uuid},
		"Email": &types.AttributeValueMemberS{Value: email},
	}
	activation := map[string]types.AttributeValue{
		":IsActive": &types.AttributeValueMemberBOOL{Value: true},
	}

	updateData := &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(tableName),
		UpdateExpression:          aws.String("set IsActive = :IsActive"),
		ExpressionAttributeValues: activation,
	}

	_, err := svc.UpdateItem(ctx, updateData)
	if err != nil {
		return false, errors.New(err.Error())
	}
	return true, nil

}
