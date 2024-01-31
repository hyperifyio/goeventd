// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package eventtrigger

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type CommandTrigger struct {
	CommandTemplate string
}

func (s *CommandTrigger) TriggerEvent(serviceName string) bool {

	if len(serviceName) == 0 {
		log.Printf("Invalid service name provided: '%s'\n", serviceName)
		return false
	}

	commandStr := fmt.Sprintf(s.CommandTemplate, serviceName)
	parts := strings.Fields(commandStr)

	if len(parts) == 0 {
		log.Printf("Invalid command template provided: '%s'\n", s.CommandTemplate)
		return false
	}

	log.Printf("Triggering service '%s' with command: '%s'\n", serviceName, commandStr)
	cmd := exec.Command(parts[0], parts[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Printf("Error triggering service '%s' using command '%s': %s\n", serviceName, commandStr, err)
		return false
	}
	log.Printf("Successfully triggered service '%s' with command '%s'\n", serviceName, commandStr)
	return true
}
