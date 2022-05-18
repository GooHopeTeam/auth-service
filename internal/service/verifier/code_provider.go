package verifier

import (
	"math/rand"
	"strconv"
	"time"
)

type codeProvider interface {
	generateCode() string
}

type codeProviderImpl struct {
	codeProvider
}

func newCodeProvider() codeProvider {
	rand.Seed(time.Now().UnixNano())
	return codeProviderImpl{}
}

func (p codeProviderImpl) generateCode() string {
	return strconv.Itoa(rand.Intn(1000000))
}
