gcl:
	@/snap/bin/go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$$(which golangci-lint) custom

lint:
	bin/custom-gcl run ./...

auth:
	@/snap/bin/go run cmd/auth/main.go

app:
	@/snap/bin/go run cmd/app/main.go

test:
	mkdir -p data
	@/snap/bin/go test -v -covermode=atomic -coverprofile=data/coverage.out ./...
	grep -v "mock" data/coverage.out > data/coverage.out.tmp
	@/snap/bin/go tool cover -html data/coverage.out.tmp -o data/coverage.html
	open data/coverage.html

auth-up:
	cd deployments/auth && docker-compose up -d

auth-down:
	cd deployments/auth && docker-compose down

.PHONY: gcl lint auth app test auth-up auth-down