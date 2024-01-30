.PHONY: build clean download

GO_SOURCES := ./messaging/messaging.go \
              ./natsclient/natsclient.go \
              ./cmd/goeventd/main.go

all: build

build: goeventd

download:
	go mod download github.com/nats-io/nats.go

goeventd: $(GO_SOURCES) download
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o goeventd ./cmd/goeventd

clean:
	rm -f goeventd
