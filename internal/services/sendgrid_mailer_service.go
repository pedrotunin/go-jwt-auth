package services

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailerService struct {
	senderName  string
	senderEmail string
	apiKey      string
}

func NewSendGridMailerService(senderName, senderEmail, apiKey string) *SendGridMailerService {
	return &SendGridMailerService{
		senderName:  senderName,
		senderEmail: senderEmail,
		apiKey:      apiKey,
	}
}

func (sgms *SendGridMailerService) SendEmail(from, to, subject, plainContent, htmlContent string) error {
	if from == "" {
		from = sgms.senderEmail
	}

	fromEmail := mail.NewEmail(sgms.senderName, from)
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(fromEmail, subject, toEmail, plainContent, htmlContent)

	client := sendgrid.NewSendClient(sgms.apiKey)

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
