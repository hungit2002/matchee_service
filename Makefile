APP_NAME=matchee-services
MAIN=./cmd/api/main.go
MYSQL_DSN='root:Tranduyhung11@tcp(127.0.0.1:3306)/matchee?charset=utf8mb4&parseTime=True&loc=Local'
.PHONY: run build tidy test lint fmt migrate-up migrate-down migrate-version migrate-cli-up migrate-cli-down migrate-cli-version migrate-new seeder-bin seed docker-build docker-up docker-down docker-logs

run:
	go run $(MAIN)

build:
	go build -o bin/$(APP_NAME) $(MAIN)
migration-bin:
	go build -o bin/migrate ./cmd/migrate

seeder-bin:
	go build -o bin/seeder ./cmd/seeder

seed:
	go run ./cmd/seeder

docker-build:
	docker build -t matchee/services:latest .

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down -v

docker-logs:
	docker compose logs -f app

tidy:
	go mod tidy

test:
	go test ./...

lint:
	gofmt -l .

fmt:
	gofmt -w .

# Migrations with golang-migrate (CLI must be installed locally)
migrate-up:
	MYSQL_DSN=$(MYSQL_DSN) go run ./cmd/migrate -direction up

migrate-down:
	MYSQL_DSN=$(MYSQL_DSN) go run ./cmd/migrate -direction down -steps 1

migrate-version:
	MYSQL_DSN=$(MYSQL_DSN) go run ./cmd/migrate -direction version

# --- golang-migrate CLI (cần cài đặt binary `migrate`) ---
# Cài (macOS): brew install golang-migrate
# Cài (khác): xem docs https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
migrate-cli-up:
	migrate -database "mysql://${MYSQL_DSN}" -path ./migrations up

migrate-cli-down:
	migrate -database "mysql://${MYSQL_DSN}" -path ./migrations down 1

migrate-cli-version:
	migrate -database "mysql://${MYSQL_DSN}" -path ./migrations version

# Tạo file migration mới (CLI):
# make migrate-new name=create_table_foo
migrate-new:
	migrate create -ext sql -dir ./migrations -seq ${name}


