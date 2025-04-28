.PHONY: help
help:
	@ echo "Usage: make [target]"
	@ echo "Targets:"
	@ echo "  help      Show this help message"
	@ echo "  tidy      Run go mod tidy"
	@ echo "  build     Build the project"
	@ echo "  test-io   Build the project with mmap and traditional methods for testing I/O performance"
	@ echo "How to run program:"
	@ echo "  start ./bin/main"
	@ echo "  run ./bin/main -h # for help"

.PHONY: tidy
tidy:
	@ echo "Running go mod tidy..."
	@ go mod tidy

.PHONY: build
build: tidy
	@ echo "Building the project..."
	@ go build -o ./bin/main -tags using_mmap_io

.PHONY: test-io
test-io: tidy
	@ echo "Building the project with mmap method for testing I/O performance..."
	@ go build -o ./bin/mmap -tags using_mmap_io
	@ echo "Building the project with traditional method for testing I/O performance..."
	@ go build -o ./bin/traditional -tags using_traditional_io