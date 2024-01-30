.PHONY: build clean download tidy

GO_EVENTD_SOURCES := ./messaging/messaging.go \
              ./natsclient/natsclient.go \
              ./cmd/goeventd/main.go

GO_EVENTD_CLIENT_SOURCES := ./cmd/goeventd-client/main.go

all: build

build: goeventd goeventd-client

tidy:
	go mod tidy

download:
	go mod download github.com/nats-io/nats.go

goeventd: $(GO_EVENTD_SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o goeventd ./cmd/goeventd

goeventd-client: $(GO_EVENTD_CLIENT_SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o goeventd-client ./cmd/goeventd-client

clean:
	rm -f goeventd goeventd-client
