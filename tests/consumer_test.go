package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"devarc.vault.apps.go-secret-consumer/internal"
	"devarc.vault.apps.go-secret-consumer/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestEchoServiceGET(t *testing.T) {
	request := httptest.NewRequest("GET", "/echo", nil)
	responseRecorder := httptest.NewRecorder()

	consumer := internal.VaultSecretConsumer{}
	consumer.InitHandlers()

	consumer.GetHttpHandlerForTesting().Handle(responseRecorder, request)
	log.Printf("message: %s", responseRecorder.Body)
	assert.Equal(t, 500, responseRecorder.Code)
	assert.Equal(t, "not implemented\n", responseRecorder.Body.String())
}

func TestEchoServicePUT(t *testing.T) {
	request := httptest.NewRequest("PUT", "/echo", nil)
	responseRecorder := httptest.NewRecorder()

	consumer := internal.VaultSecretConsumer{}
	consumer.InitHandlers()
	consumer.GetHttpHandlerForTesting().Handle(responseRecorder, request)
	log.Printf("message: %s", responseRecorder.Body)
	assert.Equal(t, 500, responseRecorder.Code)
	assert.Equal(t, "not implemented\n", responseRecorder.Body.String())
}

func TestEchoServicePATCH(t *testing.T) {
	request := httptest.NewRequest("PATCH", "/echo", nil)
	responseRecorder := httptest.NewRecorder()

	consumer := internal.VaultSecretConsumer{}
	consumer.InitHandlers()
	consumer.GetHttpHandlerForTesting().Handle(responseRecorder, request)
	log.Printf("message: %s", responseRecorder.Body)
	assert.Equal(t, 500, responseRecorder.Code)
	assert.Equal(t, "invalid method: PATCH\n", responseRecorder.Body.String())
}

func TestEchoServicePOSTSuccess(t *testing.T) {
	payload := handler.RequestPayload{
		Namespace: "test",
	}

	content, err := json.Marshal(payload)
	assert.Nil(t, err)

	request := httptest.NewRequest("POST", "/echo", bytes.NewReader(content))
	responseRecorder := httptest.NewRecorder()

	consumer := internal.VaultSecretConsumer{}
	consumer.InitHandlers()
	consumer.GetHttpHandlerForTesting().Handle(responseRecorder, request)
	log.Printf("message: %s", responseRecorder.Body)
	assert.Equal(t, 200, responseRecorder.Code)
}

func TestEchoServiceResourceNotFoundGET(t *testing.T) {
	request := httptest.NewRequest("GET", "/test", nil)
	responseRecorder := httptest.NewRecorder()

	consumer := internal.VaultSecretConsumer{}
	consumer.InitHandlers()
	consumer.GetHttpHandlerForTesting().Handle(responseRecorder, request)
	log.Printf("message: %s", responseRecorder.Body)
	assert.Equal(t, 404, responseRecorder.Code)
	assert.Equal(t, "resource not found\n", responseRecorder.Body.String())
}
