VERSION = $(shell cat ./version.json | jq .version | tr -d '"')
GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_DATE = $(shell go run ./cmd/builddate.go)

HOME_OS ?= $(shell go env GOOS)
ifeq ($(HOME_OS),linux)
	TARGET_OS ?= linux
else ifeq ($(HOME_OS),darwin)
	TARGET_OS ?= darwin
else ifeq ($(HOME_OS),windows)
	TARGET_OS ?= windows
else
	$(error System $(LOCAL_OS) is not supported.)
endif

ifeq ($(HOME_OS),windows)
	TARGET_FILE ?= ./build/timeouts.exe
else
	TARGET_FILE ?= ./build/timeouts
endif

# Usage: `make build`
build:
	@echo Building project...
	go build -ldflags "-s -w -X main.version=${VERSION} -X main.commitHash=${GIT_COMMIT} -X \"main.buildDate=${BUILD_DATE}\"" -o $(TARGET_FILE)
	@echo Built project! Use "$(TARGET_FILE)" to execute the project.

# Usage: `make clean`
clean:
	@echo Cleaning...
	go clean
	rm -rf build
	@echo Done!

# Usage: `make clean`
clean:
	go clean
