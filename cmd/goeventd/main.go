// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/hyperifyio/goeventd/messaging"
    "github.com/hyperifyio/goeventd/natsclient"
    "github.com/nats-io/nats.go"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
)

type ServiceConfig struct {
    Subject     string `json:"subject"`
    ServiceName string `json:"service"`
}

type Config struct {
    NATSServer    string          `json:"natsServer"`
    ServiceConfig []ServiceConfig `json:"services"`
}

var (
    natsServer = flag.String("nats", os.Getenv("GOEVENTD_NATS_URL"), "The NATS server URL")
    subject = flag.String("subject", os.Getenv("GOEVENTD_SUBJECT"), "The NATS subject to subscribe to")
    service = flag.String("service", os.Getenv("GOEVENTD_SERVICE"), "The SystemD service to trigger")
    shutdownAfterRun = flag.Bool("once", false, "Shutdown the service after one successful SystemD service execution")
    configFilePath = flag.String("config", os.Getenv("GOEVENTD_CONFIG"), "Path to configuration file")
)

func main() {

    var msgHandler messaging.MessageHandler
    var config *Config
    var err error

    flag.Parse()

    // Get NATS URL from environment variable or use default
    if *natsServer == "" {
        defaultURL := nats.DefaultURL
        natsServer = &defaultURL
    }

    // If a config file is provided, use it
    if *configFilePath != "" {
        config, err = readConfig(*configFilePath)
        if err != nil {
            log.Fatalf("Error reading config file: %s", err)
        }
    } else {

        if *subject == "" {
            log.Fatal("--subject cannot be empty")
        }

        if *service == "" {
            log.Fatal("--service cannot be empty")
        }

        // Otherwise, use command line arguments
        config = &Config{
            NATSServer: *natsServer,
            ServiceConfig: []ServiceConfig{
                {
                    Subject:     *subject,
                    ServiceName: *service,
                },
            },
        }
    }

    // Initialize NATS client
    msgHandler = &natsclient.NATSClient{}
    err = msgHandler.Initialize(messaging.Config{ServerURL: *natsServer})
    if err != nil {
        fmt.Println("Error initializing NATS client:", err)
        return
    }
    defer msgHandler.Close()

    // Subscribe to a subject
    for _, sc := range config.ServiceConfig {
        err := msgHandler.Subscribe(sc.Subject, func(msg string) {
            fmt.Printf("Event(%s): %s\n", sc.Subject, msg)
            triggerSystemDService(sc.ServiceName)
        })
        if err != nil {
            log.Fatalf("Error subscribing to subject %s: %s", sc.Subject, err)
        } else {
            fmt.Printf("Successfully subscribed to: %s\n", sc.Subject)
        }
    }

    // Keep the program running
    select {}
}

func triggerSystemDService(service string) bool {
    fmt.Printf("Executing: systemctl start %s\n", service)
    cmd := exec.Command("systemctl", "start", service)
    err := cmd.Run()
    if err != nil {
        fmt.Printf("Error triggering SystemD service %s using 'systemctl start %s': %s\n", service, service, err)
        return false
    }
    fmt.Printf("Successfully triggered SystemD service: %s\n", service)
    return true
}

func readConfig(filePath string) (*Config, error) {
    file, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    var config Config
    err = json.Unmarshal(file, &config)
    if err != nil {
    return nil, err
    }
    return &config, nil
}
