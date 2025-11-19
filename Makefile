.PHONY: run-dev build-backend up down

PROJECT_NAME=

run-dev-backend:
	cd backend && go run cmd/api/main.go

run-dev-frontend:
	cd frontend && yarn dev

up:
	docker-compose up -d --build

down:
	docker-compose down

mock-gen:
	# Команда для генерации моков (позже добавить mockery)
	echo "Generating mocks..."
