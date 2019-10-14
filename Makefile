INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
TEST ?= $(shell go list ./... | grep -v -e vendor)

BINARY:=pouta
GO ?= GO111MODULE=on $(SYSTEM) go

dbcreate: ## Crate user and Create database
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Setup database$(RESET)"
	mysql -uroot -h127.0.0.1 -P3306 -e 'create database if not exists pouta DEFAULT CHARACTER SET ujis;'

dbdrop: ## Drop database
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Drop database$(RESET)"
	mysql -uroot -h127.0.0.1 -P3306 -e 'drop database pouta;'

dbmigrate: depends dbcreate
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Migrate database$(RESET)"
	sql-migrate up pouta 
depends:
	GO111MODULE=off go get -v github.com/rubenv/sql-migrate/...

run:
	$(GO) run main.go --config ./misc/develop.toml server

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)
