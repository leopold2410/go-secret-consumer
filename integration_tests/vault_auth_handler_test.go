package integrationtests

import (
	"context"
	"testing"

	"devarc.vault.apps.go-secret-consumer/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestVaultAuthHandlerGetVaultToken(t *testing.T) {
	vaultAddress := "http://localhost:8200"
	appRoleId := "8b44fc80-276d-4001-2f2b-a5f39a9caeb1"
	appSecretId := "6a52480a-f938-65e3-2fb0-bd050b689986"
	client, err := handler.CreateVaultClient(vaultAddress)
	assert.Nil(t, err)
	authHandler := handler.VaultAuthHandler{}
	authHandler.Initialize(context.Background(), client)
	err = authHandler.Login(appRoleId, appSecretId)
	assert.Nil(t, err)
}

func TestVaultAuthHandlerGetKV2Secret(t *testing.T) {
	ctx := context.Background()
	vaultAddress := "http://localhost:8200"
	appRoleId := "8b44fc80-276d-4001-2f2b-a5f39a9caeb1"
	appSecretId := "615d1581-1780-59c6-de0a-7de10757a655"
	client, err := handler.CreateVaultClient(vaultAddress)
	assert.Nil(t, err)
	authHandler := handler.VaultAuthHandler{}
	authHandler.Initialize(ctx, client)
	err = authHandler.Login(appRoleId, appSecretId)
	assert.Nil(t, err)
	secretEnginePath := "kv"
	secretPath := "test"
	kvEngin := client.KVv2(secretEnginePath)
	assert.NotNil(t, kvEngin)
	secret, err := kvEngin.Get(ctx, secretPath)
	assert.Nil(t, err)
	metadata := secret.VersionMetadata
	if metadata != nil {
		assert.NotEmpty(t, handler.GetTimeAsString(metadata.CreatedTime))
	}
}
