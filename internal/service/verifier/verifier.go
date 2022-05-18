package verifier

type Verifier interface {
	Send(email string, data map[string]string) error
	Check(email, code string) (map[string]string, error)
}
