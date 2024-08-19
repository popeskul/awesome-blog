# swagger
SWAGGER_UI_VERSION := v5.4.2
SWAGGER_UI_DIR := ./static/swagger-ui

.PHONY: generate-swagger
generate-swagger: generate-swagger-api install-swagger-ui
	@echo "Generating Swagger documentation..."
	@cp ./api/openapi/api.yaml $(SWAGGER_UI_DIR)/swagger.yaml
	@awk '{gsub("https://petstore.swagger.io/v2/swagger.json", "./swagger.yaml")}1' $(SWAGGER_UI_DIR)/swagger-initializer.js > $(SWAGGER_UI_DIR)/swagger-initializer.js.tmp && mv $(SWAGGER_UI_DIR)/swagger-initializer.js.tmp $(SWAGGER_UI_DIR)/swagger-initializer.js
	@echo "Swagger documentation generated successfully"


.PHONY: generate-swagger-api
generate-swagger-api: .ensure-dir
	@echo "Generating API code..."
	@which oapi-codegen > /dev/null || go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	@oapi-codegen -config api/openapi/cfg.yaml api/openapi/api.yaml > gen/api/openapi.gen.go
	@echo "API code generated successfully."

.PHONY: install-swagger-ui
install-swagger-ui:
	@echo "Installing Swagger UI..."
	@mkdir -p $(SWAGGER_UI_DIR)
	@curl -sSL https://github.com/swagger-api/swagger-ui/archive/$(SWAGGER_UI_VERSION).tar.gz | tar -xz -C /tmp
	@cp -R /tmp/swagger-ui-$(SWAGGER_UI_VERSION:v%=%)/dist/* $(SWAGGER_UI_DIR)
	@rm -rf /tmp/swagger-ui-$(SWAGGER_UI_VERSION:v%=%)
	@echo "Swagger UI installed successfully"

.PHONY: .ensure-dir
.ensure-dir:
	@mkdir -p gen/api/

.PHONY: test
test:
	@go test ./...

.PHONY: lint
lint:
	@golangci-lint run

# Migrate
MIGRATIONS_DIR = ./migrations
DB_DSN = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: install-migrate
install-migrate:
	@echo "Installing migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: migrate-setup
migrate-setup: install-migrate migrate-up

# make migrate-create
.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

# make migrate-up DB_USER=user DB_PASSWORD=password DB_HOST=localhost DB_PORT=5432 DB_NAME=blogdb DB_SSLMODE=disable
.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up

# make migrate-down DB_USER=user DB_PASSWORD=password DB_HOST=localhost DB_PORT=5432 DB_NAME=blogdb DB_SSLMODE=disable
.PHONY: migrate-down
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down

.PHONY: migrate-down-to
migrate-down-to:
	@read -p "Enter number of migrations to roll back: " num; \
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down $$num

.PHONY: migrate-status
migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" version

.PHONY: migrate-force
migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" force $$version

.PHONY: migrate-drop
migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" drop

.PHONY: run
run:
	go run cmd/main.go
