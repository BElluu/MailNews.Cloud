package services

import (
	"MailNews.Subscriber/common"
	"MailNews.Subscriber/database"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
)

func SendActivateEmail(toEmail string) {
	client := common.DynamoDBSession()
	subscriber := database.GetSubscriber2(toEmail, client)

	body := "There is your activate link:" + subscriber.ActivateURL
	var recip = []*string{&subscriber.Email}
	SendEmailSES(body, "MailNews.Cloud - Activate newsletter.", "xxx", recip)
}

func SendEmailSES(messageBody string, subject string, fromEmail string, recipient []*string) {

	session := common.AmazonSESSesion()

	svc := ses.New(session)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: recipient,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(messageBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(fromEmail),
	}
	_, err := svc.SendEmail(input)
	if err != nil {
		log.Println("Error sending mail - ", err)
		return
	}
	log.Println("Email sent successfully.")
}
