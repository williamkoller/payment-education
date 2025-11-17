package email

import (
	"fmt"

	"github.com/resend/resend-go/v3"
)

type ResendClient struct {
	sdkClient *resend.Client
	from      string
}

func NewResendClient(apiKey, fromAddress string) *ResendClient {
	client := resend.NewClient(apiKey)
	return &ResendClient{
		sdkClient: client,
		from:      fromAddress,
	}
}

func (r *ResendClient) SendEmail(to, subject, html string) error {
	params := &resend.SendEmailRequest{
		From:    r.from,
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}

	resp, err := r.sdkClient.Emails.Send(params)

	if err != nil {
		return fmt.Errorf("resend email send failed: %w", err)
	}

	fmt.Printf("Resend email sent: id=%s\n", resp.Id)

	return nil
}
