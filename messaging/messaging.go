// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package messaging

type Config struct {
	ServerURL string
}

type MessageHandler interface {
	Initialize(config Config) error
	Subscribe(subject string, action func(msg string)) error
	Close()
}
