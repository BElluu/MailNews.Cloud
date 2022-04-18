package database

import (
	"MailNews.Subscriber/models"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"log"
)

const FeedTable = "MailNewsFeeds"

func AddFeed(ctx context.Context, feedItem models.FeedItem, client *dynamodb.Client) {
	svc := client
	tableName := FeedTable
	id := uuid.New().String()
	feedMap := map[string]types.AttributeValue{
		"UUID":        &types.AttributeValueMemberS{Value: id},
		"Title":       &types.AttributeValueMemberS{Value: feedItem.Title},
		"Link":        &types.AttributeValueMemberS{Value: feedItem.Link},
		"PublishDate": &types.AttributeValueMemberS{Value: feedItem.PublishDate.Format("02-01-2006 15:01:05")},
		"Source":      &types.AttributeValueMemberS{Value: feedItem.Source},
	}

	input := &dynamodb.PutItemInput{
		Item:      feedMap,
		TableName: aws.String(tableName),
	}

	_, err := svc.PutItem(ctx, input)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func GetConfigValue(ctx context.Context, configProperty string, client *dynamodb.Client) (map[string]types.AttributeValue, error) {
	svc := client
	tableName := "MailNewsConfig"

	configProp := map[string]types.AttributeValue{
		"Name": &types.AttributeValueMemberS{Value: configProperty},
	}

	getPropertyValue := &dynamodb.GetItemInput{
		Key:       configProp,
		TableName: aws.String(tableName),
	}

	configResult, err := svc.GetItem(ctx, getPropertyValue)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	if configResult.Item == nil {
		return nil, nil
	}
	return configResult.Item, nil
}
