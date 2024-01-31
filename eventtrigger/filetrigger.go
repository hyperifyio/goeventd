// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package eventtrigger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type FileTrigger struct {
	FileFormat string
}

func (f *FileTrigger) TriggerEvent(serviceName string) bool {

	filePath := fmt.Sprintf(f.FileFormat, serviceName)
	currentTime := time.Now().Local()

	// Create the file if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening or creating service '%s' file '%s': %s\n", serviceName, filePath, err)
		return false
	}

	if err := file.Close(); err != nil {
		log.Printf("Error closing service '%s' file '%s': %s\n", serviceName, filePath, err)
		return false
	}

	// Update file's access and modification times
	err = os.Chtimes(filePath, currentTime, currentTime)
	if err != nil {
		log.Printf("Error updating times for service '%s' and file '%s': %s\n", serviceName, filePath, err)
		return false
	}

	log.Printf("Successfully triggered service '%s' by file modification to: %s\n", serviceName, filePath)
	return true
}
