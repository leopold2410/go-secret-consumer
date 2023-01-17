package internal

import "devarc.vault.apps.go-secret-consumer/internal/handler"

type VaultSecretConsumer struct {
	httpHandler *handler.HttpHandler
}

func (srv *VaultSecretConsumer) GetHttpHandlerForTesting() *handler.HttpHandler {
	return srv.httpHandler
}

func (srv *VaultSecretConsumer) InitHandlers() {
	vaultHandler := handler.VaultSecretHandler{}
	vaultHandler.Initialize()

	echo := handler.EchoHandler{}

	httpHandler := handler.HttpHandler{}
	httpHandler.Init()
	httpHandler.RegisterHandler("/fetch", &vaultHandler)
	httpHandler.RegisterHandler("/echo", &echo)

	srv.httpHandler = &httpHandler
}
