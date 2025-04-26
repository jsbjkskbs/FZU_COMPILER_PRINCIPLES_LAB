.PHONY: help
help:
	@ echo "Usage: make [target]"
	@ echo ""
	@ echo "Targets:"
	@ echo "  help        Show this help message"
	@ echo "  tidy        Run go mod tidy"
	@ echo "  lexer       Run the lexer"
	@ echo "  parser      Run the parser"

.PHONY: tidy
tidy:
	@ go mod tidy

.PHONY: lexer
lexer: tidy
	@ go run main.go lexer

.PHONY: parser
parser: tidy
	@ go run main.go parser
