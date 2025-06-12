### Settings & utilities ###

.PHONY: gen lint test coverage run prep check

GREEN  := \033[1;32m
RESET  := \033[0m

define banner
	@echo ""; echo "$(GREEN)$(1)$(RESET)"
endef

# A single place for the env you want when you run the server
SERVER_ENV = \
	POSTGRES_DB=tablename \
	POSTGRES_HOST=localhost \
	POSTGRES_PASSWORD=asdf \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=postgres


### targets ###

gen:
	$(call banner,Generating templ.go files…)
	@templ generate

lint:
	$(call banner,Installing & running GolangCI-Lint…)
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
	@golangci-lint run ./...

test:
	$(call banner,Running tests…)
	@go test -race -v -coverprofile=coverage.out ./...

coverage:
	$(call banner,Checking code coverage…)
	@sh .scripts/check_coverage.sh

### Composed targets ###

# Local dev run: build templates, then start server
run: gen
	$(call banner,Setting up environment and starting webserver…)
	@env $(SERVER_ENV) go run main.go -l

# “prep” for commit: everything should pass before you push
prep: gen lint test coverage

# CI lint only
check: lint
