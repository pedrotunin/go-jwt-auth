package services

type MailerService interface {
	SendEmail(from, to, subject, plainContent, htmlContent string) error
}
