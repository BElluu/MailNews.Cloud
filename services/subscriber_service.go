package services

import (
	"MailNews.Subscriber/common"
	"MailNews.Subscriber/common/validator"
	"MailNews.Subscriber/database"
	"MailNews.Subscriber/models"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"time"
)

func Subscribe(email string) error {
	client := common.CreateLocalClient()

	if !Email.IsValid(email) {
		return errors.New("email address is not valid")
	}

	if database.GetSubscriber2(email, client).Email != "" {
		return errors.New("email exists")
	}
	newSubscriber(email, client)
	return nil

}

func UnSubscribe(uuid string, email string) error {
	client := common.CreateLocalClient()
	if exist, _ := isSubscriberExist(uuid, email); exist {
		err := database.DeleteSubscriber(uuid, email, client)
		if err != nil {
			return errors.New(err.Error())
		}
		return nil
	}
	return errors.New("email does not exist")
}

func ActivateSubscription(uuid string, email string) error {
	client := common.CreateLocalClient()
	isActive := database.GetSubscriber2(email, client).IsActive
	if isActive == true {
		return errors.New("subscription is activated")
	}
	if exist, _ := isSubscriberExist(uuid, email); exist {
		_, err := database.ActiveSubscriber(uuid, email, client)
		if err != nil {
			return errors.New(err.Error())
		}
		return nil
	}
	return errors.New("email does not exist")

}

func isSubscriberExist(uuid string, email string) (bool, error) {
	client := common.CreateLocalClient()
	subscriber, err := database.GetSubscriber(uuid, email, client)
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
		CreatedDate:    time.Now().UTC().Format("02-01-2006 15:01:05"),
	}
	database.AddSubscriber(subscriber, client)
}
