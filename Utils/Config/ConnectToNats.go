package config

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

var NatsConnection *nats.Conn

func ConnectToNATS() {
	natsConnection, natsErr := nats.Connect("tls://connect.ngs.global:4222", nats.UserCredentials("./NatsCred.creds"))

	if natsErr != nil {
		log.Fatal("Failed to connect to nats", natsErr)
	}

	//assign this to global variable
	NatsConnection = natsConnection

	fmt.Println("NATS Connected Successfully")

}
