package database

import (
	"MailNews.Subscriber/models"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetConfigValue(configProperty string, client *dynamodb.Client) (*models.Config, error) {
	svc := client
	tableName := "MailNewsConfig"

	configProp := map[string]types.AttributeValue{
		"Name": &types.AttributeValueMemberS{Value: configProperty},
	}

	getPropertyValue := &dynamodb.GetItemInput{
		Key:       configProp,
		TableName: aws.String(tableName),
	}

	configResult, err := svc.GetItem(context.Background(), getPropertyValue)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	config := models.Config{}
	err = attributevalue.UnmarshalMap(configResult.Item, &config)
	return &config, nil
}
