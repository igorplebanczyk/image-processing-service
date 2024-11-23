package emails

import (
	"fmt"
	"github.com/wneessen/go-mail"
)

type Service struct {
	client *mail.Client
	sender string
}

func NewService(sender, password string) (*Service, error) {
	client, err := mail.NewClient(
		"smtp.gmail.com",
		mail.WithPort(587),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithUsername(sender),
		mail.WithPassword(password),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create email client: %v", err)
	}

	return &Service{
		sender: sender,
		client: client,
	}, nil
}

func (s *Service) SendText(to []string, subject, body string) error {
	message := mail.NewMsg()

	err := message.From(s.sender)
	if err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	for _, recipient := range to {
		err = message.To(recipient)
		if err != nil {
			return fmt.Errorf("failed to set recipient: %v", err)
		}
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	err = s.client.DialAndSend(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
