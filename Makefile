gcl:
	@/snap/bin/go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$$(which golangci-lint) custom

lint:
	bin/custom-gcl run ./...

server:
	@/snap/bin/go run cmd/$(service)/main.go

test:
	mkdir -p data
	@/snap/bin/go test -v -covermode=atomic -coverprofile=data/coverage.out ./...
	grep -v "mock" data/coverage.out > data/coverage.out.tmp
	@/snap/bin/go tool cover -html data/coverage.out.tmp -o data/coverage.html
	open data/coverage.html

up:
	cd deployments/$(service) && docker-compose up -d

down:
	cd deployments/$(service) && docker-compose down

migrate:
	migrate create -ext sql -dir migrations/$(service) -seq -digits 2 $(name)

.PHONY: gcl lint server test up down migrate