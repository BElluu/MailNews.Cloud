package sender

import (
	"MailNews.Cloud/backend/common"
	dbservice "MailNews.Cloud/database/services"
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"html/template"
	"log"
)

type Recipient struct {
	ToEmails []string
	CcEmails []string
}

type templateData struct {
	ActivateLink string
}

func SendActivateEmail(toEmail string) {
	client := common.CreateLocalClient()
	subscriber := dbservice.GetSubscriber2(toEmail, client)
	// TODO Change Path
	tmpl := template.Must(template.ParseFiles("/home/bartek/Programming/MailNews.Cloud/sender/templates/activate.html"))

	data := templateData{ActivateLink: subscriber.ActivateURL}
	var tpl bytes.Buffer
	tmpl.Execute(&tpl, data)
	result := tpl.String()
	println(result)

	body := result
	var recip = []*string{&subscriber.Email}
	err := SendEmailSES(body, "MailNews.Cloud - Activate newsletter.", "xxx", recip)
	if err != nil {
		err := dbservice.DeleteSubscriber(subscriber.UUID, subscriber.Email, client)
		if err != nil {
			return // PANIC - log it!
		}
	}
}

func SendEmailSES(messageBody string, subject string, fromEmail string, recipient []*string) error {

	session := common.AmazonSESSession()

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
		return err
	}
	log.Println("Email sent successfully.")
	return nil
}
