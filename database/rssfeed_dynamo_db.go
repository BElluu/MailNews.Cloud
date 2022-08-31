package database

import (
	"MailNews.Subscriber/models"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"log"
)

const FeedTable = "MailNewsFeeds"
const AWSTable = "AwsNews"
const AzureTable = "AzureNews"
const GCPTable = "GoogleCloudNews"

func AddFeed(feedItem models.FeedItem, client *dynamodb.Client, provider string) {
	svc := client
	tableName := ""
	switch provider {
	case "Aws":
		tableName = AWSTable
	case "Azure":
		tableName = AzureTable
	case "Google":
		tableName = GCPTable
	}

	id := uuid.New().String()
	feedMap := map[string]types.AttributeValue{
		"UUID":        &types.AttributeValueMemberS{Value: id},
		"Title":       &types.AttributeValueMemberS{Value: feedItem.Title},
		"Link":        &types.AttributeValueMemberS{Value: feedItem.Link},
		"PublishDate": &types.AttributeValueMemberS{Value: feedItem.PublishDate.Format("02-01-2006 15:01:05")},
		//"Provider":    &types.AttributeValueMemberS{Value: feedItem.Provider},
		"Sent": &types.AttributeValueMemberBOOL{Value: feedItem.Sent},
	}

	input := &dynamodb.PutItemInput{
		Item:      feedMap,
		TableName: aws.String(tableName),
	}

	_, err := svc.PutItem(context.Background(), input)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func FeedFromProviderExist(provider string, client *dynamodb.Client) bool {

	svc := client
	tableName := FeedTable
	filter := expression.Name("Provider").Equal(expression.Value(provider)).And(expression.Name("Sent").Equal(expression.Value(false)))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		panic(err)
	}

	out, err := svc.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		panic(err)
	}

	if len(out.Items) == 0 {
		return false
	}
	return true
}

func GetFeedsToSend(provider string, client *dynamodb.Client) []models.FeedItem {
	svc := client
	tableName := SubscriberTable
	filter := expression.Name("Provider").Equal(expression.Value(provider)).And(expression.Name("Sent").Equal(expression.Value(false)))
	proj := expression.NamesList(expression.Name("UUID"), expression.Name("Title"),
		expression.Name("Description"),
		expression.Name("Link"))
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

	var feeds []models.FeedItem

	for _, value := range out.Items {
		item := models.FeedItem{}
		err = attributevalue.UnmarshalMap(value, &feeds)
		if err != nil {
			println("wtf")
		}
		feeds = append(feeds, item)
	}
	return feeds
}
