package database

import (
	"MailNews.Subscriber/models"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

func AddSubscriber(subscriber models.Subscriber, client *dynamodb.Client) {
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

	_, err := svc.PutItem(context.Background(), input)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Println("Added " + subscriber.Email + " to table" + tableName)
}

func DeleteSubscriber(uuid string, email string, client *dynamodb.Client) error {
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

	_, err := svc.DeleteItem(context.Background(), input)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func ActiveSubscriber(uuid string, email string, client *dynamodb.Client) (bool, error) {
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

	_, err := svc.UpdateItem(context.Background(), updateData)
	if err != nil {
		return false, errors.New(err.Error())
	}
	return true, nil
}

func GetSubscriber(uuid string, email string, client *dynamodb.Client) (*dynamodb.GetItemOutput, error) {
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

	subscriberResult, err := svc.GetItem(context.Background(), getSubscriber)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	if subscriberResult.Item == nil {
		return nil, nil
	}
	return subscriberResult, nil
}

func GetSubscribers(activeSubscribers bool, client *dynamodb.Client) []models.Subscriber {
	svc := client
	tableName := SubscriberTable
	filter := expression.Name("IsActive").Equal(expression.Value(activeSubscribers))
	proj := expression.NamesList(expression.Name("Email"),
		expression.Name("ActivateURL"),
		expression.Name("UnSubscribeURL"))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(proj).Build()
	if err != nil {
		panic(err)
	}

	out, err := svc.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		panic(err)
	}

	var subscribers []models.Subscriber

	for _, value := range out.Items {
		item := models.Subscriber{}
		err = attributevalue.UnmarshalMap(value, &item)
		if err != nil {
			println("wtf")
		}
		subscribers = append(subscribers, item)
	}

	if len(subscribers) == 0 {
		return nil
	}
	return subscribers
}

func GetSubscriber2(email string, client *dynamodb.Client) models.Subscriber {
	svc := client
	tableName := SubscriberTable
	filter := expression.Name("Email").Equal(expression.Value(email))
	proj := expression.NamesList(expression.Name("Email"),
		expression.Name("ActivateURL"),
		expression.Name("UnSubscribeURL"),
		expression.Name("IsActive"))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(proj).Build()
	if err != nil {
		panic(err)
	}

	out, err := svc.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		panic(err)
	}

	var subscriber models.Subscriber
	for _, value := range out.Items {
		err = attributevalue.UnmarshalMap(value, &subscriber)
		if err != nil {
			println("wtf")
		}
	}
	return subscriber
}
