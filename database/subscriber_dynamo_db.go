package database

import (
	"MailNews.Subscriber/models"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

const SubscriberTable = "Subscribers"

func AddSubscriber(ctx context.Context, subscriber models.Subscriber, client *dynamodb.Client) {
	svc := client
	tableName := SubscriberTable

	subscriberMap := map[string]types.AttributeValue{
		"UUID":           &types.AttributeValueMemberS{Value: subscriber.UUID},
		"Email":          &types.AttributeValueMemberS{Value: subscriber.Email},
		"ActivateURL":    &types.AttributeValueMemberS{Value: subscriber.ActivateURL},
		"UnSubscribeURL": &types.AttributeValueMemberS{Value: subscriber.UnSubscribeURL},
		"IsActive":       &types.AttributeValueMemberBOOL{Value: subscriber.IsActive},
		"CreatedDate":    &types.AttributeValueMemberS{Value: subscriber.CreatedDate},
	}

	input := &dynamodb.PutItemInput{
		Item:      subscriberMap,
		TableName: aws.String(tableName),
	}

	_, err := svc.PutItem(ctx, input)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Println("Added " + subscriber.Email + " to table" + tableName)
}

func DeleteSubscriber(ctx context.Context, uuid string, email string, client *dynamodb.Client) (bool, error) {
	svc := client
	tableName := SubscriberTable

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

func ActiveSubscriber(ctx context.Context, uuid string, email string, client *dynamodb.Client) (bool, error) {
	svc := client
	tableName := SubscriberTable

	subscriber := map[string]types.AttributeValue{
		"UUID":  &types.AttributeValueMemberS{Value: uuid},
		"Email": &types.AttributeValueMemberS{Value: email},
	}
	activation := map[string]types.AttributeValue{
		":IsActive": &types.AttributeValueMemberBOOL{Value: true},
	}

	updateData := &dynamodb.UpdateItemInput{
		Key:                       subscriber,
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

func GetSubscriber(ctx context.Context, uuid string, email string, client *dynamodb.Client) (*dynamodb.GetItemOutput, error) {
	svc := client
	tableName := SubscriberTable

	subscriber := map[string]types.AttributeValue{
		"UUID":  &types.AttributeValueMemberS{Value: uuid},
		"Email": &types.AttributeValueMemberS{Value: email},
	}

	getSubscriber := &dynamodb.GetItemInput{
		Key:       subscriber,
		TableName: aws.String(tableName),
	}

	subscriberResult, err := svc.GetItem(ctx, getSubscriber)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	if subscriberResult.Item == nil {
		return nil, nil
	}
	return subscriberResult, nil
}
