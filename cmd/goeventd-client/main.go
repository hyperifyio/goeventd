// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package main

import (
	"flag"
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

var (
	natsServer = flag.String("nats", os.Getenv("NATS_URL"), "The NATS server URL")
	subject = flag.String("subject", "", "The NATS subject to subscribe to")
	message = flag.String("message", "", "The NATS message")
)

func main() {

	flag.Parse()

	// Get NATS URL from environment variable or use default
	if *natsServer == "" {
		defaultURL := nats.DefaultURL
		natsServer = &defaultURL
	}

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

	log.Printf("Published message '%s' to subject '%s'", *message, *subject)

}
