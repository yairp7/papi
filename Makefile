GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BUILD_DIR=build

run:
	$(GOCMD) run .

build:
	mkdir -p $(BUILD_DIR)/
	$(GOCMD) build -ldflags="-s -w" -gcflags=all="-l" -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	rm -fr ./$(BUILD_DIR)

watch:
	$(eval PACKAGE_NAME=$(shell head -n 1 go.mod | cut -d ' ' -f2))
	docker run -it --rm -w /go/src/$(PACKAGE_NAME) -v $(shell pwd):/go/src/$(PACKAGE_NAME) -p ${PORT}:${PORT} cosmtrek/air

test: 
	$(GOTEST) -v -race ./... $(OUTPUT_OPTIONS)

gen-config:
	echo "CONFIG=$(shell base64 -i config.json)" > .env