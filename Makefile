run-dev:
	air -c .air.toml

run:
	./bin/link-service

generate-config:
	cp .env.example .env