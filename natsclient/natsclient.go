package natsclient

import (
	"github.com/hyperifyio/goeventd/messaging"
	"github.com/nats-io/nats.go"
)

type NATSClient struct {
	conn *nats.Conn
}

func (c *NATSClient) Initialize(config messaging.Config) error {
	var err error
	c.conn, err = nats.Connect(config.ServerURL)
	return err
}

func (c *NATSClient) Subscribe(subject string, action func(msg string)) error {
	_, err := c.conn.Subscribe(subject, func(m *nats.Msg) {
		action(string(m.Data))
	})
	return err
}

func (c *NATSClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
