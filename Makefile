.PHONY: help
help:
	@ echo "Usage: make [target]"
	@ echo "Targets:"
	@ echo "  help      Show this help message"
	@ echo "  tidy      Run go mod tidy"
	@ echo "  build     Build the project"
	@ echo "  test-io   Build the project with mmap and traditional methods for testing I/O performance"
	@ echo "  fmt       Format the code and imports"
	@ echo "How to run program:"
	@ echo "  start ./bin/main"
	@ echo "  run ./bin/main -h # for help"

.PHONY: tidy
tidy:
	@ echo "Running go mod tidy..."
	@ go mod tidy

.PHONY: build
build: tidy fmt
	@ echo "Building the project..."
	@ go build -o ./bin/main -tags using_mmap_io

.PHONY: test-io
test-io: tidy fmt
	@ echo "Building the project with mmap method for testing I/O performance..."
	@ go build -o ./bin/mmap -tags using_mmap_io
	@ echo "Building the project with traditional method for testing I/O performance..."
	@ go build -o ./bin/traditional -tags using_traditional_io

.PHONY: fmt
fmt:
	@ echo "Checking if goimports is installed..."
	@ command -v goimports >/dev/null 2>&1 || { echo >&2 "goimports is not installed. Please install it using 'go install golang.org/x/tools/cmd/goimports@latest'."; exit 1; }
	@ echo "Checking if gofumpt is installed..."
	@ command -v gofumpt >/dev/null 2>&1 || { echo >&2 "gofumpt is not installed. Please install it using 'go install mvdan.cc/gofumpt@latest'."; exit 1; }
	@ echo "Formatting imports..."
	@ goimports -w -local app .
	@ echo "Formatting the code..."
	@ gofumpt -l -w .