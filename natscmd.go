package main

import (
	"flag"
	"fmt"
	"github.com/goombaio/namegenerator"
	nats "github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

const (
	pubCmd = "pub"
	subCmd = "sub"
	reqCmd = "req"
	repCmd = "rep"
)

// Version of this tool
var Version string = "dev"

func main() {
	natsURL := flag.String("nats", "nats://localhost:4222", "NATS server URL")
	cmd := flag.String("cmd", "sub", "sub, pub or req")
	subject := flag.String("subject", "test", "NATS subject to use")
	message := flag.String("message", "Hello World", "Message to send")
	timeout := flag.Duration("timeout", time.Second*30, "Subscriber timeout")
	certsDir := flag.String("certs", "", "(Optional) Path to directory with client.pem, client-key.pem and ca.pem")
	creds := flag.String("creds", "", "(Optional) Path to credentials file")
	v := flag.Bool("v", false, "Version")

	flag.Parse()

	if *v {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	name := randomName()
	log.Printf("Name: %s\n", name)
	// Connect Options.
	opts := []nats.Option{
		nats.Name(name),
		nats.NoEcho(),
	}

	if *certsDir != "" {
		opts = append(opts, nats.RootCAs(fmt.Sprintf("%s/ca.pem", *certsDir)))
		opts = append(opts, nats.ClientCert(fmt.Sprintf("%s/client.pem", *certsDir), fmt.Sprintf("%s/client-key.pem", *certsDir)))
	}

	if *creds != "" {
		opts = append(opts, nats.UserCredentials(*creds))
	}

	// Connect to a server
	nc, err := nats.Connect(*natsURL, opts...)
	if err != nil {
		log.Fatal(err)
	}

	switch *cmd {
	case pubCmd:
		err := nc.Publish(*subject, []byte(*message))
		if err != nil {
			log.Fatal(err)
		}
	case subCmd:
		// Simple Sync Subscriber
		sub, err := nc.SubscribeSync(*subject)
		m, err := sub.NextMsg(*timeout)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Received message '%s' from '%s'", string(m.Data), m.Subject)
	case reqCmd:
		msg, err := nc.Request(*subject, []byte(name+":"+*message), 2*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(string(msg.Data))
	case repCmd:
		sub, err := nc.SubscribeSync(*subject)
		if err != nil {
			log.Fatal(err)
		}

		msg, err := sub.NextMsg(*timeout)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(string(msg.Data))
		msg.Respond([]byte(name + ":" + *message))
	}
	nc.Close()
}

func randomName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	return nameGenerator.Generate()
}
