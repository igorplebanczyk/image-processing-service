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

func (s *Service) SendOTP(recipient, subject, issuer, code string) error {
	return s.SendHTML(recipient, subject, template(issuer, code))
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

func template(issuer, code string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OTP Verification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            width: 100%%;
            background-color: #ffffff;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        h1 {
            color: #333333;
            font-size: 24px;
            margin: 0 0 20px 0;
        }
        .otp {
            font-size: 36px;
            font-weight: bold;
            color: #007bff;
            background-color: #ffffff;
            padding: 15px 30px;
            border-radius: 8px;
            display: inline-block;
            margin: 20px 0 40px 0;
        }
        .footer {
            margin-top: 30px;
            font-size: 14px;
            color: #888888;
        }
        table {
            border-spacing: 0;
            width: 100%%;
        }
    </style>
</head>
<body padding: 20px;">
    <table align="center" role="presentation" cellpadding="0" cellspacing="0" width="100%%">
        <tr>
            <td align="center">
                <table class="container" cellpadding="0" cellspacing="0" role="presentation">
                    <tr>
                        <td align="center">
                            <h1>%s</h1>
                        </td>
                    </tr>
                    <tr>
                        <td align="center">
                            <div class="otp">%s</div>
                        </td>
                    </tr>
                    <tr>
                        <td align="center" class="footer">
                            If you did not request this code, please ignore this email.
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, issuer, code)
}
