package main

import (
	"flag"
	"github.com/nats-io/nats.go"
	"log"
)

var (
	natsServer = flag.String("nats", "nats://localhost:4222", "The NATS server URL")
	subject = flag.String("subject", "", "The NATS subject to subscribe to")
	message = flag.String("message", "", "The NATS message")
)

func main() {

	flag.Parse()

	if *subject == "" {
		log.Fatal("Subject cannot be empty")
	}

	// Connect to the NATS server.
	nc, err := nats.Connect(*natsServer)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Publish a message.
	err = nc.Publish(*subject, []byte(*message))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Published message '%s' to subject '%s'", message, subject)

}
