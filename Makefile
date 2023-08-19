SOURCES := $(shell find . -mindepth 2 -name "main.go")
DESTS := $(patsubst %.go,%,$(SOURCES))
ALL := main $(DESTS)

all: $(ALL)
	@echo $@: Building Targets $^

main:
	@echo $@: Building main
	go build -buildvcs -o dist/$@ $@.go

$(DESTS):
	@echo $@: Building $@ to ${shell dirname dist/$@}
	go build -buildvcs -o ${shell dirname dist/$@} $@.go

dep:
	go mod tidy

clean:
	go clean
	rm -f $(ALL)

.PHONY: clean

run:
	arelo -t . -p '**/*.go' -i '**/.*' -i '**/*_test.go' -- go run .