.PHONY: build clean

GO_SOURCES := ./messaging/messaging.go \
              ./natsclient/natsclient.go \
              ./cmd/goeventd/main.go

all: build

build: goeventd

goeventd: $(GO_SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o goeventd ./cmd/goeventd

clean:
	rm -f goeventd
