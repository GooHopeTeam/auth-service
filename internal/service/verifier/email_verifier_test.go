package verifier

import (
	"reflect"
	"testing"
)
import "github.com/stretchr/testify/assert"

func getVerifier() EmailVerifier {
	return EmailVerifier{newMailerMock(), newStorageMock(), newCodeProviderMock()}
}

func TestEmailVerifier_Send(t *testing.T) {
	assert := assert.New(t)
	data := map[string]string{"message": "hello world"}
	ver := getVerifier()
	email := "some_email@gmail.com"
	err := ver.Send(email, data)
	assert.NoError(err)
	savedData, _ := ver.storage.Get(email)
	message, containsMessage := savedData["message"]
	assert.True(containsMessage)
	assert.Equal(data["message"], message)
}

func TestEmailVerifier_Check_Unset(t *testing.T) {
	assert := assert.New(t)
	ver := getVerifier()
	email := "some_email@gmail.com"
	data, err := ver.Check(email, "0")
	assert.NoError(err)
	assert.Nil(data)
}

func TestEmailVerifier_Check_Sent(t *testing.T) {
	assert := assert.New(t)
	data := map[string]string{"message": "hello world"}
	ver := getVerifier()
	email := "some_email@gmail.com"
	err := ver.Send(email, data)
	assert.NoError(err)
	savedData, err := ver.Check(email, "1")
	assert.NoError(err)
	assert.NotNil(data)
	assert.True(reflect.DeepEqual(data, savedData))
}
