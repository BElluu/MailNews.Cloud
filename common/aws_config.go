package common

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
)

const Region = "eu-central-1"
const AccessKeyID = "dummy"
const SecretAccessKey = ""
const SessionToken = ""

func DynamoDBSession() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	return client
}

func AmazonSESSesion() *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(Region)})
	if err != nil {
		log.Println("Error occurred while creating aws session", err)
	}
	return sess
}
