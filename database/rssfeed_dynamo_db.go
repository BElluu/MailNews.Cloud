package database

import (
	"MailNews.Subscriber/models"
	"context"
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
