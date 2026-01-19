run-dev:
	air -c .air.toml

run:
	./bin/link-service

generate-config:
	cp .env.example .env


test-cov:
	 go test -coverprofile=coverage.xml ./...

test:
	 go test  ./...

build:
	go build -o ./bin/app