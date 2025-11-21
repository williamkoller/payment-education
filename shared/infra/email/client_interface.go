package email

type EmailClient interface {
	SendEmail(to, subject, html string) error
}
