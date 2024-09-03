GOOS = linux
GOARCH = arm
ARM_VERSION=7
PROJECT_NAME = authn-svc
MAIN_FILE = cmd/server/main.go
BUILD_PATH = build
SHELL := /bin/bash


init-build-dirs:
	@rm -rf $(BUILD_PATH)
	@mkdir $(BUILD_PATH)
	$(info $(BUILD_PATH) was created)

build-api: init-build-dirs
	$(info Will build API for $(GOOS) and $(GOARCH))
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(ARM_VERSION) go build -o $(BUILD_PATH)/$(PROJECT_NAME) $(MAIN_FILE)
	@echo "API build completed"

clean:
	@rm -rf $(BUILD_PATH)

run:
	go run $(MAIN_FILE)

lint:
	$(HOME)/go/bin/golangci-lint run
