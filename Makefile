GIT_COMMIT ?= noversion
GIT_COMMIT_SHORT ?= noversion

LDFLAGS = "-X main.Version=$(GIT_COMMIT)"

objects = cmd/query/main

all: $(objects)

$(objects):
	go build -ldflags $(LDFLAGS) $@.go

dep:
	go mod tidy

vendor:
	go mod vendor

clean:
	go clean
	rm -f $(objects)

.PHONY: clean
