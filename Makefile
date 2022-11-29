build:
	@CGO_ENABLED=0 go build -ldflags "-s -w" -o "dist/app" github.com/WeNeedThePoh/nba-stats-api/cmd/server

run:
	go run cmd/server/main.go

run_docker:
	docker-compose up --build -d

test:
	go test ./...

lint: tools.golangci-lint
	@CGO_ENABLED=0 ./bin/golangci-lint-latest run -n

# ----------- #
# -- TOOLS -- #
# ----------- #
tools.golangci-lint:
	@command -v ./bin/golangci-lint-latest >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "$(OK_COLOR)==> installing golangci-lint version latest $(NO_COLOR)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest; \
		mv ./bin/golangci-lint ./bin/golangci-lint-latest; \
	fi
