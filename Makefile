include .env
export

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down:
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume:
	docker volume rm pg-data
.PHONY: docker-rm-volume

migrate-create:
	migrate create -ext sql -dir migrations 'marketplace'
.PHONY: migrate-create

migrate-up:
	migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' up
.PHONY: migrate-up

migrate-down:
	echo "y" | migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' down
.PHONY: migrate-down

test:
	go test -v ./...

cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
.PHONY: coverage-html

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out
.PHONY: coverage

mockgen:
	mockgen github.com/amiosamu/vk-internship/internal/service Service > internal/mocks/servicemocks/service_mock.go
	mockgen github.com/amiosamu/vk-internship/internal/repo Repo > internal/mocks/repomocks/repo_mock.go
.PHONY: mockgen

swag:
	swag init -g internal/app/app.go --parseInternal --parseDependency