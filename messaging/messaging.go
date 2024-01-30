package messaging

type Config struct {
	ServerURL string
}

type MessageHandler interface {
	Initialize(config Config) error
	Subscribe(subject string, action func(msg string)) error
	Close()
}
