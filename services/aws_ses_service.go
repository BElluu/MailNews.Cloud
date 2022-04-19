package services

import (
	"MailNews.Subscriber/common"
	"MailNews.Subscriber/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
)

func SendEmailSES(messageBody string, subject string, fromEmail string, recipient models.Recipient) {

	session := common.AmazonSESSesion()
	var recipients []*string

	for _, r := range recipient.ToEmails { // in result from dynamo
		recipient := r
		recipients = append(recipients, &recipient)
	}

	svc := ses.New(session)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: recipients,
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
	log.Println("Email sent successfully to: ", recipient.ToEmails)
}
