package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"log"
	"time"
)

type News struct {
	UUID        string
	Title       string
	Link        string
	PublishDate *time.Time
	Provider    string
	Sent        bool
}

func AddNews(news News, client *dynamodb.Client, provider string) {
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
		"Title":       &types.AttributeValueMemberS{Value: news.Title},
		"Link":        &types.AttributeValueMemberS{Value: news.Link},
		"PublishDate": &types.AttributeValueMemberS{Value: news.PublishDate.Format("02-01-2006 15:01:05")},
		"Sent":        &types.AttributeValueMemberBOOL{Value: news.Sent},
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

func NewsFromProviderExist(provider string, client *dynamodb.Client) bool {
	//TODO I do not remember I need this method. Check it!
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

func GetNewsToSend(provider string, client *dynamodb.Client) []News {
	svc := client
	tableName := AWSTable // TODO check all tables
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

	var news []News

	for _, value := range out.Items {
		item := News{}
		err = attributevalue.UnmarshalMap(value, &news)
		if err != nil {
			println("wtf")
		}
		news = append(news, item)
	}
	return news
}
