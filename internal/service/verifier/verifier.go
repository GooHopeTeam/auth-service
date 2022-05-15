package verifier

type Verifier[Data any] interface {
	Send(email string, data *Data) error
	Check(email, code string) (*Data, error)
}
