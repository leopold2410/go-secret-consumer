package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
)

type VaultSecretHandler struct {
	client      *vault.Client
	authHandler *VaultAuthHandler
	ctx         context.Context
}

type RequestPayload struct {
	AppRole      string `json:"appRole" description:"appRole for login"`
	AppSecret    string `json:"appSecret" description:"appSecret for login"`
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
	AppRole string `json:"appRole"`
}

type Message struct {
	Secret SecretInfo `json:"secret"`
	Token  TokenInfo  `json:"token"`
}

func CreateVaultClient(vaultAddress string) (*vault.Client, error) {
	config := vault.DefaultConfig()
	config.Address = vaultAddress
	return vault.NewClient(config)
}

func (srv *VaultSecretHandler) Initialize() error {
	srv.ctx = context.Background()
	authHandler := VaultAuthHandler{}
	vaultClient, err := CreateVaultClient(os.Getenv("VAULT_ADDR"))
	if err != nil {
		return err
	}
	authHandler.Initialize(srv.ctx, vaultClient)
	srv.authHandler = &authHandler
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

	if payload.AppRole == "" {
		http.Error(w, fmt.Sprintf("invalid Parameter AppRole: %s", payload.AppRole), http.StatusBadRequest)
		return
	}

	if payload.AppSecret == "" {
		http.Error(w, fmt.Sprintf("invalid Parameter AppSecret: %s", payload.AppSecret), http.StatusBadRequest)
		return
	}

	err = srv.authHandler.Login(payload.AppRole, payload.AppSecret)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	namespace := payload.Namespace
	if namespace != "" {
		srv.client.SetNamespace(namespace)
	}
	secretEnginePath := payload.SecretEngine

	secretEngine := srv.client.KVv2(secretEnginePath)
	if secretEngine == nil {
		http.Error(w, fmt.Sprintf("SecretEngine: %s not found", payload.SecretEngine), http.StatusBadRequest)
		return
	}

	secretPath := payload.SecretPath
	kvSecret, err := secretEngine.Get(srv.ctx, secretPath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	metadata := kvSecret.VersionMetadata
	if metadata == nil {
		http.Error(w, "no metadata found", http.StatusBadRequest)
		return
	}

	secretInfo := SecretInfo{
		Version:      metadata.Version,
		SecretPath:   secretPath,
		EnginePath:   secretEnginePath,
		CreationTime: GetTimeAsString(metadata.CreatedTime),
		DeletionTime: GetTimeAsString(metadata.DeletionTime),
		LeaseId:      kvSecret.Raw.LeaseID,
		RequestId:    kvSecret.Raw.RequestID,
	}

	tokenInfo := TokenInfo{AppRole: payload.AppRole}
	message := Message{Token: tokenInfo, Secret: secretInfo}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func GetTimeAsString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.String()
}
