package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"devarc.vault.apps.go-secret-consumer/internal"
)

func main() {
	log.Println("Consumer started")
	port := flag.Int("port", 5050, "listen port")
	if port == nil {
		log.Fatal("port not specified")
	}

	consumer := internal.VaultSecretConsumer{}
	consumer.InitHandlers()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
