package services

import (
	"MailNews.Subscriber/common/validator"
	"MailNews.Subscriber/database"
	"MailNews.Subscriber/models"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

func Subscribe(email string) (bool, error) {
	if Email.IsValid(email) {
		client := database.CreateLocalClient()
		newSubscriber(email, client)
		return true, nil
	}
	return false, errors.New("email address is not valid")
}

func UnSubscribe(uuid string, email string) (bool, error) {
	client := database.CreateLocalClient()
	if exist, _ := isSubscriberExist(uuid, email); exist {
		_, err := database.DeleteSubscriber(context.TODO(), uuid, email, client)
		if err != nil {
			return false, errors.New(err.Error())
		}
		return true, nil
	}
	return false, errors.New("email does not exist")
}

func ActivateSubscription(uuid string, email string) (bool, error) {
	client := database.CreateLocalClient()
	if exist, _ := isSubscriberExist(uuid, email); exist {
		_, err := database.ActiveSubscriber(context.TODO(), uuid, email, client)
		if err != nil {
			return false, errors.New(err.Error())
		}
		return true, nil
	}
	return false, errors.New("email does not exist")

}

func isSubscriberExist(uuid string, email string) (bool, error) {
	client := database.CreateLocalClient()
	subscriber, err := database.GetSubscriber(context.TODO(), uuid, email, client)
	if err != nil {
		return false, errors.New(err.Error())
	}
	if subscriber == nil {
		return false, nil
	}
	return true, nil
}

func newSubscriber(email string, client *dynamodb.Client) {
	id := uuid.New().String()
	var subscriber = models.Subscriber{
		UUID:           id,
		Email:          email,
		ActivateURL:    "http://localhost:8080/activate/?email=" + email + "&uuid=" + id,
		UnSubscribeURL: "http://localhost:8080/unsubscribe/?email=" + email + "&uuid=" + id,
		IsActive:       false,
	}
	database.AddSubscriber(context.TODO(), subscriber, client)
}
