package mailer

type Mailer interface {
	Send(email, message string) error
}
