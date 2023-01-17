package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	vault "github.com/hashicorp/vault/api"
	vault_approle "github.com/hashicorp/vault/api/auth/approle"
)

type AuthContext struct {
	AppRole string
}

type VaultSecretHandler struct {
	authContext *AuthContext
	client      *vault.Client
	authHandler *vault_approle.AppRoleAuth
	ctx         context.Context
}

type RequestPayload struct {
	Namespace    string `json:"namespace" description:"vault namespace (enterprise feature)"`
	SecretEngine string `json:"secretEngine" description:"path to secret engine in vault"`
	SecretPath   string `json:"secretPath" description:"path to secret"`
}

type SecretInfo struct {
	EnginePath   string `json:"enginePath"`
	SecretPath   string `json:"secretPath"`
	Version      int    `json:"version"`
	CreationTime string `json:"creationTime"`
	DeletionTime string `json:"deletionTime"`
	LeaseId      string `json:"leaseId"`
	RequestId    string `json:"requestId"`
}

type TokenInfo struct {
	TokenId       string `json:"tokenId"`
	RequestId     string `json:"requestId"`
	LeaseId       string `json:"leaseId"`
	LeaseDuration int    `json:"leaseDuration"`
	Ttl           int    `json:"ttl"`
	AppRole       string `json:"appRole"`
}

type Message struct {
	Token  TokenInfo  `json:"token"`
	Secret SecretInfo `json:"secret"`
}

func createTokenInfo(authInfo *vault.Secret, appRole string) (*TokenInfo, error) {
	tokenId, err := authInfo.TokenID()
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		TokenId:       tokenId,
		RequestId:     authInfo.RequestID,
		LeaseId:       authInfo.LeaseID,
		LeaseDuration: authInfo.LeaseDuration,
		Ttl:           authInfo.WrapInfo.TTL,
		AppRole:       appRole,
	}, nil
}

func createAuthHandler(approle string, secret string) (*vault_approle.AppRoleAuth, error) {
	secret_id := vault_approle.SecretID{
		FromString: secret,
	}

	return vault_approle.NewAppRoleAuth(approle, &secret_id)
}

func createVaultClient() (*vault.Client, error) {
	config := vault.DefaultConfig()
	config.ReadEnvironment()
	return vault.NewClient(config)
}

func (srv *VaultSecretHandler) Initialize() error {
	appRole := os.Getenv("VAULT_APP_ROLE")
	appRoleSecret := os.Getenv("VAULT_APP_ROLE_SECRET")
	authHandler, err := createAuthHandler(appRole, appRoleSecret)
	if err != nil {
		return err
	}
	vaultClient, err := createVaultClient()
	if err != nil {
		return err
	}
	srv.authContext = &AuthContext{AppRole: appRole}
	srv.authHandler = authHandler
	srv.client = vaultClient
	return nil
}

/*
/fetch
*/
func (srv *VaultSecretHandler) Get(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "not implemented", 500)
}

func (srv *VaultSecretHandler) Put(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "not implemented", 500)
}

func ParseRequestPayload(req *http.Request) (*RequestPayload, error) {
	var payload RequestPayload
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

func (srv *VaultSecretHandler) Post(w http.ResponseWriter, req *http.Request) {
	payload, err := ParseRequestPayload(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authSecret, err := srv.authHandler.Login(srv.ctx, srv.client)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	tokenInfo, err := createTokenInfo(authSecret, srv.authContext.AppRole)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	namespace := payload.Namespace
	if namespace != "" {
		srv.client.SetNamespace(namespace)
	}
	secretEnginePath := payload.SecretEngine
	secretPath := payload.SecretPath
	kvSecret, err := srv.client.KVv2(secretEnginePath).Get(srv.ctx, secretPath)

	secretInfo := SecretInfo{
		Version:      kvSecret.VersionMetadata.Version,
		SecretPath:   secretPath,
		EnginePath:   secretEnginePath,
		CreationTime: kvSecret.VersionMetadata.CreatedTime.String(),
		DeletionTime: kvSecret.VersionMetadata.DeletionTime.String(),
		LeaseId:      kvSecret.Raw.LeaseID,
		RequestId:    kvSecret.Raw.RequestID,
	}

	message := Message{Token: *tokenInfo, Secret: secretInfo}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}
