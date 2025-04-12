.PHONY: build run clean test docker-build docker-up docker-down test-endpoints

# Development commands
build:
	go build -o bin/main .

run:
	go run main.go

test:
	go test ./...

clean:
	rm -rf bin/

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Test commands
test-endpoints:
	@echo "Testing health endpoint..."
	@curl -s http://localhost:8080/health | jq
	@echo "\nTesting database connection..."
	@curl -s http://localhost:8080/ping-db | jq

# Database commands
db-migrate:
	docker-compose exec mysql mysql -u root -prootpassword portal < init.sql

# Helper commands
ps:
	docker-compose ps

restart:
	docker-compose restart 