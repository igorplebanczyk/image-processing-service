package emails

import (
	"fmt"
	"github.com/wneessen/go-mail"
	"log/slog"
)

type Service struct {
	client *mail.Client
	sender string
}

func NewService(host, sender, password string) (*Service, error) {
	client, err := mail.NewClient(
		host,
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

func (s *Service) SendText(recipient, subject, body string) error {
	return s.send(recipient, subject, body, mail.TypeTextPlain)
}

func (s *Service) SendHTML(recipient, subject, body string) error {
	return s.send(recipient, subject, body, mail.TypeTextHTML)
}

func (s *Service) send(recipient, subject, body string, contentType mail.ContentType) error {
	message := mail.NewMsg()

	err := message.From(s.sender)
	if err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	err = message.To(recipient)
	if err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	message.Subject(subject)
	message.SetBodyString(contentType, body)

	go func() {
		err = s.client.DialAndSend(message)
		if err != nil {
			slog.Error("Failed to send email", "error", err)
		}
	}()

	return nil
}
