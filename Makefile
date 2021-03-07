.PHONY: all
all: help


.PHONY: build
build: ## Build project and put output binary in /bin folder
	@go build -o bin/wcrawler cmd/wcrawler/main.go


.PHONY: test
test: ## Run the tests of the project
	@go test -race ./...


.PHONY: coverage
coverage: ## Run the tests of the project and print out coverage
	@go test -cover ./...


.PHONY: coverage-report
coverage-report: ## Run the tests of the project and show line by line coverage in the browser
	@go test -coverprofile=coverage.txt ./...
	@go tool cover -html=coverage.txt


.PHONY: bench
bench: ## Run benchmarks of the project
	@go test -benchmem -bench . -run=XXX ./...


.PHONY: lint
lint: ## Run linters
	@gofmt -l .
	@go vet ./...


.PHONY: clean
clean: ## Remove temporary and build related files
	@rm -f coverage.txt
	@rm -f bin/*


.PHONY: docs-server
docs-server: ## Start godoc server to allow navigation of documentation
	@echo "Documentation @ http://127.0.0.1:6060"
	@godoc -http=:6060


.PHONY: find_todo
find_todo: ## Find all TODOs in the code
	@grep --color=always --include=\*.go -PnRe '(//|/*).*TODO' ./ || true


.PHONY: count
count: ## Count number of lines in go files
	@echo "Lines of code:"
	@find . -type f -name "*.go" | xargs wc -l


GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)


.PHONY: help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${YELLOW}%-16s${GREEN}%s${RESET}\n", $$1, $$2}' $(MAKEFILE_LIST)
