package mailer

import (
	"net/smtp"
)

type MailerImpl struct {
	Mailer
	sender   string
	smtpHost string
	smtpPort string
	smtpAuth smtp.Auth
}

func New(smtpHost, smtpPort, sender, password string) Mailer {
	return &MailerImpl{smtpHost: smtpHost, smtpPort: smtpPort, sender: sender,
		smtpAuth: smtp.PlainAuth("", sender, password, smtpHost)}
}

func (m MailerImpl) Send(email, message string) error {
	return smtp.SendMail(m.smtpHost+":"+m.smtpPort, m.smtpAuth, m.sender, []string{email}, []byte(message))
}
