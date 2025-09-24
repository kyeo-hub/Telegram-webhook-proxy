# Makefile for Telegram Bot API Proxy

BINARY=telegram-proxy
GOBUILD=go build
GOCLEAN=go clean
GOINSTALL=go install
GOTEST=go test

# Detect OS
ifeq ($(OS),Windows_NT)
    TARGET := $(BINARY).exe
    RM := del /Q
else
    TARGET := $(BINARY)
    RM := rm -f
endif

all: build

build:
	$(GOBUILD) -o $(TARGET) .

clean:
	$(GOCLEAN)
	$(RM) $(TARGET)

test:
	$(GOTEST) -v

install:
	$(GOINSTALL)

run: build
	./$(TARGET)

.PHONY: all build clean test install run