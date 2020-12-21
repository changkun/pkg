package mailgun

import (
	"context"

	"github.com/mailgun/mailgun-go/v4"
)

// SendEmail sends an email via mailgun
func SendEmail(ctx context.Context, domain, apikey, endpoint string,
	sender, subject, body, recipient string) error {
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(domain, apikey)
	mg.SetAPIBase(endpoint)

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	// Send the message with a 10 second timeout
	_, _, err := mg.Send(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
