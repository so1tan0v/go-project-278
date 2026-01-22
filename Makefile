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

docker-build:
	docker build -t go-project-278 .

docker-run:
	docker run -d --rm -p 8080:8080 --env-file ./.env go-project-278