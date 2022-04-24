package common

import (
	"context"
	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	aws1 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
)

const Region = "eu-central-1"
const AccessKeyID = "dummy"
const SecretAccessKey = ""
const SessionToken = ""

func CreateLocalClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("local"),
		config.WithEndpointResolverWithOptions(aws2.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws2.Endpoint, error) {
				return aws2.Endpoint{URL: "http://127.0.0.1:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws2.Credentials{AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "MailNews",
			},
		}),
	)
	if err != nil {
		panic(err)
	}
	return dynamodb.NewFromConfig(cfg)
}

func DynamoDBSession() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	return client
}

func AmazonSESSesion() *session.Session {
	sess, err := session.NewSession(&aws1.Config{Region: aws1.String(Region)})
	if err != nil {
		log.Println("Error occurred while creating aws session", err)
	}
	return sess
}
