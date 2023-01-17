package main

import (
	"flag"
	"log"

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
	consumer.StartServer(*port)
}
