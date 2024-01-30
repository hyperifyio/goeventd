
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/hyperifyio/goeventd/messaging"
    "github.com/hyperifyio/goeventd/natsclient"
    "io/ioutil"
    "log"
    "os/exec"
)

type ServiceConfig struct {
    Subject     string `json:"subject"`
    ServiceName string `json:"serviceName"`
}

type Config struct {
    NATSServer    string          `json:"natsServer"`
    ServiceConfig []ServiceConfig `json:"services"`
}

var (
    natsServer = flag.String("nats", "nats://localhost:4222", "The NATS server URL")
    subject = flag.String("subject", "", "The NATS subject to subscribe to")
    serviceName = flag.String("service", "", "The SystemD service to trigger")
    shutdownAfterRun = flag.Bool("once", false, "Shutdown the service after one successful SystemD service execution")
    configFilePath = flag.String("config", "", "Path to configuration file")
)

func main() {

    var msgHandler messaging.MessageHandler

    flag.Parse()

    var config *Config
    var err error

    // If a config file is provided, use it
    if *configFilePath != "" {
        config, err = readConfig(*configFilePath)
        if err != nil {
            log.Fatalf("Error reading config file: %s", err)
        }
    } else {
        // Otherwise, use command line arguments
        config = &Config{
            NATSServer: *natsServer,
            ServiceConfig: []ServiceConfig{
                {
                    Subject:     *subject,
                    ServiceName: *serviceName,
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
            triggerSystemDService(sc.ServiceName)
        })
        if err != nil {
            log.Fatalf("Error subscribing to subject %s: %s", sc.Subject, err)
        }
    }

    // Keep the program running
    select {}
}

func triggerSystemDService(serviceName string) bool {
    cmd := exec.Command("systemctl", "start", serviceName)
    err := cmd.Run()
    if err != nil {
        fmt.Printf("Error triggering SystemD service %s: %s\n", serviceName, err)
        return false
    }
    fmt.Printf("Successfully triggered SystemD service: %s\n", serviceName)
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
