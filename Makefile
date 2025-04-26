.PHONY: help
help:
	@ echo "Usage: make [target]"
	@ echo "Targets:"
	@ echo "  help      Show this help message"
	@ echo "  tidy      Run go mod tidy"
	@ echo "  build     Build the project with mmap and traditional version"
	@ echo "How to run program:"
	@ echo "  start ./bin/mmap or ./bin/traditional"
	@ echo "  ./bin/mmap -h or ./bin/traditional -h for help"

.PHONY: tidy
tidy:
	@ echo "Running go mod tidy..."
	@ go mod tidy

.PHONY: build
build: tidy
	@ echo "Building the project with mmap method..."
	@ go build -o ./bin/mmap -tags using_mmap_io
	@ echo "Building the project with traditional method..."
	@ go build -o ./bin/traditional -tags using_traditional_io