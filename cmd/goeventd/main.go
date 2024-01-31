// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/hyperifyio/goeventd/eventtrigger"
    "github.com/hyperifyio/goeventd/messaging"
    "github.com/hyperifyio/goeventd/natsclient"
    "github.com/nats-io/nats.go"
    "log"
    "os"
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
    natsServer = flag.String("nats", getEnvOrDefault("GOEVENTD_NATS_URL", nats.DefaultURL), "The NATS server URL")
    subject = flag.String("subject", getEnvOrDefault("GOEVENTD_SUBJECT", ""), "The NATS subject to subscribe to")
    service = flag.String("service", getEnvOrDefault("GOEVENTD_SERVICE", ""), "The SystemD service to trigger")
    shutdownAfterRun = flag.Bool("once", false, "Shutdown the service after one service triggering")
    configFilePath = flag.String("config", getEnvOrDefault("GOEVENTD_CONFIG", ""), "Path to configuration file")
    eventType = flag.String("event-type", getEnvOrDefault("GOEVENTD_EVENT_TYPE","command"), "The event triggering mechanism (command/file)")
    fileTemplate = flag.String(getEnvOrDefault("GOEVENTD_FILE_TEMPLATE","file-template"), "./%s.trigger", "The path to trigger services by file modifications. The service name can be provided in %s placeholder. (when event-type is file)")
    commandTemplate = flag.String(getEnvOrDefault("GOEVENTD_COMMAND_TEMPLATE","command-template"), "systemctl start %s", "The command to trigger services. The service name can be provided in %s placeholder. (when event-type is command)")
)

func main() {

    var trigger eventtrigger.EventTrigger
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

    // Initialize appropriate trigger
    switch *eventType {

    case "command":
        if *commandTemplate == "" {
            log.Fatal("--command-template cannot be empty")
        }
        trigger = &eventtrigger.CommandTrigger{CommandTemplate: *commandTemplate}

    case "file":
        if *fileTemplate == "" {
            log.Fatal("--file-template cannot be empty")
        }
        trigger = &eventtrigger.FileTrigger{FileFormat: *fileTemplate}

    default:
        log.Fatalf("Unknown event type: %s", *eventType)
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
            log.Printf("Event '%s' received with '%s'\n", sc.Subject, msg)
            if err := trigger.TriggerEvent(sc.ServiceName); !err {
                log.Printf("Error triggering event for service %s\n", sc.ServiceName)
            }
            if *shutdownAfterRun {
                log.Printf("Closing service since --once was enabled.\n")
                os.Exit(0)
            }
        })
        if err != nil {
            log.Fatalf("Error subscribing to subject %s: %s", sc.Subject, err)
        } else {
            log.Printf("Successfully subscribed to: %s\n", sc.Subject)
        }
    }

    // Keep the program running
    select {}
}

func readConfig(filePath string) (*Config, error) {
    file, err := os.ReadFile(filePath)
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

func getEnvOrDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}