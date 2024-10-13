GIT_COMMIT_SHORT=$(shell git rev-parse --short HEAD)
COMMIT_TIME=$(shell git show -s --format=%cd --date=iso-strict HEAD)

LDFLAGS=-X github.com/tendant/chi-demo/app.GitCommit=$(GIT_COMMIT_SHORT) \
        -X github.com/tendant/chi-demo/app.CommitTime=$(COMMIT_TIME) \


SOURCES := $(shell find . -mindepth 2 -name "main.go")
DESTS := $(patsubst ./%/main.go,dist/%,$(SOURCES))
ALL := dist/main $(DESTS)

all: $(ALL)
	@echo $@: Building Targets $^

dist/main:
ifneq (,$(wildcard main.go))
	$(echo Bulding main.go)
	go build -buildvcs -ldflags "$(LDFLAGS)" -o $@ main.go
endif

#dist/main:
#	@echo Building $^ into $@
#	test -f main.go && go build -buildvcs -o $@ $^

dist/%: %/main.go
	@echo $@: Building $^ to $@
	go build -buildvcs -ldflags "$(LDFLAGS)" -o $@ $^

dep:
	go mod tidy

clean:
	go clean
	rm -f $(ALL)

.PHONY: clean

run:
	arelo -t . -p '**/*.go' -i '**/.*' -i '**/*_test.go' -- go run .