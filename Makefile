include .envrc


# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: welcome
welcome:
	@printf "\033[0m\n"
	@printf "\033[33m  ██████╗██╗  ██╗███████╗███████╗███████╗\n"
	@printf "\033[33m ██╔════╝██║  ██║██╔════╝██╔════╝██╔════╝\n"
	@printf "\033[33m ██║     ███████║█████╗  ███████╗███████╗\n"
	@printf "\033[33m ██║     ██╔══██║██╔══╝  ╚════██║╚════██║\n"
	@printf "\033[33m ╚██████╗██║  ██║███████╗███████║███████║\n"
	@printf "\033[33m  ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝\n"
	@printf "\033[0m\n"


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api: welcome
	@go run ./cmd/api -cs=${GREEN_DB_CS} -smtp-username ${SMTP_USER} -smtp-password ${SMTP_PW} -smtp-host ${SMTP_HOST} -smtp-sender ${SMTP_SENDER}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	@psql ${GREEN_DB_CS}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	@migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	@migrate -path ./migrations -database ${GREEN_DB_CS} up


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...