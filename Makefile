# Include variables from the .envrc file
ifneq (,$(wildcard ./.envrc))
    include .envrc
endif

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## copy-hooks: adds script to run before git push
copy-hooks:
	chmod +x scripts/hooks/*
	cp -r scripts/hooks .git/.

## expvar: add environment variable required for testing
expvar:
	eval $(cat .envrc)

## generate: generate code from specs
generate:
	go generate ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: run the api application
run:
	go run . -db-dsn=${TLTV_DB_DSN}

## db/psql: connect to the database using psql
db/psql:
	psql ${TLTV_DB_DSN}

## db/migrations/new name=$1: create a new database migration
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./db/migrations ${name}

## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./db/migrations -database ${TLTV_DB_DSN} up

## db/migrations/down: apply all down database migrations
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path ./db/migrations -database ${TLTV_DB_DSN} down

## db/dump: pg_dump current tltv database
db/dump: confirm
	pg_dump --dbname=${TLTV_DB_DSN} -F t >> db/testdata/tltv_db_$(shell date +%s).tar

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running tests...'

## audit/pipeline: tidy dependencies and format, vet and test all code (race on)
audit/pipeline:
	make audit
	go test -race -vet=off ./... -coverprofile=coverage.out
## audit/local: tidy dependencies and format, vet and test all code (race off)
audit/local:
	make audit
	go test -vet=off ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o cover.html

## staticcheck:  detect bugs, suggest code simplifications, and point out dead code
staticcheck:
	staticcheck ./...

## coverage
coverage:
	go tool cover -func coverage.out \
	| grep "total:" | awk '{print ((int($$3) > 80) != 1) }'

## coverage report
report:
	go test -vet=off ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o cover.html

install-golang-ci:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

ci-lint: install-golang-ci
	golangci-lint run

install-govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest

vuln: install-govulncheck
	govulncheck ./*.go

key:
	openssl ecparam -name prime256v1 -genkey -noout -out ecprivatekey.pem

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date +"%Y-%m-%dT%H:%M:%S%Z")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build: build the cmd/api application
build:
	@echo 'Building api...'
	rm ecprivatekey.pem
	go build -ldflags=${linker_flags} -o=./bin/api ./api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./api