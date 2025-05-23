gcl:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$$(which golangci-lint) custom

lint:
	bin/custom-gcl run ./...

auth:
	go run cmd/auth/main.go

app:
	go run cmd/app/main.go

test:
	mkdir -p data
	go test -v -covermode=atomic -coverprofile=data/coverage.out ./...
	grep -v "mock" data/coverage.out > data/coverage.out.tmp
	go tool cover -html data/coverage.out.tmp -o data/coverage.html
	open data/coverage.html

app-up:
	cd deployments/app && docker-compose up -d

app-down:
	cd deployments/app && docker-compose down

auth-up:
	cd deployments/auth && docker-compose up -d

auth-down:
	cd deployments/auth && docker-compose down

.PHONY: gcl lint auth app test auth-up auth-down app-up app-down