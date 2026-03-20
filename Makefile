#!make
GOCMD=go
SQLC_CMD=sqlc
GOOSECMD=goose
GOBUILDFLAGS=-ldflags=''
BUILD_DIR=out
EXE_NAME=transActor

.PHONY: all build run test clean dbclean tidy fmt vet up down reset
all: generate up tidy test build

build: $(BUILD_DIR)
	@$(GOCMD) build -o $(BUILD_DIR)/$(EXE_NAME)

$(BUILD_DIR):
	@echo "Creating out directory..."
	@mkdir -p $(BUILD_DIR)

run: build
	@echo "running..."
	@$(GOCMD) run .

tidy: fmt vet
	@echo "Tidying go modules..."
	@$(GOCMD) mod tidy

fmt:
	@echo "Formatting files..."
	@$(GOCMD) fmt ./...

vet:
	@echo "Vetting project..."
	@$(GOCMD) vet ./...

clean:
	@echo "Tidying project..."
	@rm -rf $(BUILD_DIR)/*
	@rm -f $(EXE_NAME)

dbclean:
	@echo "wiping db..."
	@rm -rf $(DB_DIR)
	@mkdir -p $(DB_DIR)
	@touch $(GOOSE_DBSTRING)

$(GOOSE_DBSTRING):
	@echo "Creating database at $(GOOSE_DBSTRING)"
	@mkdir -p $(DB_DIR)
	@touch $(GOOSE_DBSTRING)

up: $(GOOSE_DBSTRING)
	@echo "Migrating Database..."
	@$(GOOSECMD) up

down: $(GOOSE_DBSTRING)
	@echo "Resetting Database..."
	@$(GOOSECMD) down

reset: $(GOOSE_DBSTRING)
	@echo "Resetting Database..."
	@$(GOOSECMD) reset

validate: $(GOOSE_DBSTRING)
	@echo "Validating migrations..."
	@$(GOOSECMD) validate

generate:
	@echo "Generating go internal files for sql..."
	@sqlc generate

