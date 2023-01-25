package handler

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
	vault_approle "github.com/hashicorp/vault/api/auth/approle"
)

type AuthMethod int

const (
	AppRole AuthMethod = 0
)

type VaultAuthHandler struct {
	authMethod         AuthMethod
	appRoleAuthHandler *vault_approle.AppRoleAuth
	vaultClient        *vault.Client
	ctx                context.Context
}

func (srv *VaultAuthHandler) Initialize(ctx context.Context, client *vault.Client) {
	srv.vaultClient = client
	srv.ctx = ctx
}

func (srv *VaultAuthHandler) Login(appRole string, appSecret string) error {
	secret, err := srv.loginInternal(appRole, appSecret)
	if err != nil {
		return err
	}
	srv.vaultClient.SetToken(secret.Auth.ClientToken)
	return nil
}

func (srv *VaultAuthHandler) loginInternal(appRole string, appSecret string) (*vault.Secret, error) {
	if srv.authMethod == AppRole {
		authHandler, err := createAuthHandler(appRole, appSecret)
		if err != nil {
			return nil, err
		}
		return authHandler.Login(srv.ctx, srv.vaultClient)
	}
	return nil, fmt.Errorf("invalid auth method")
}

func createAuthHandler(approle string, secret string) (*vault_approle.AppRoleAuth, error) {
	secret_id := vault_approle.SecretID{
		FromString: secret,
	}

	return vault_approle.NewAppRoleAuth(approle, &secret_id)
}
