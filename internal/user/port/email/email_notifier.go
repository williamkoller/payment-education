package port_email_notifier 

type EmailNotifier interface {
	SendWelcomeEmail(name, email string) error
}