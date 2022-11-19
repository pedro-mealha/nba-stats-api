MIGRATE_VERSION=v4.15.1

# Determine platform
PLATFORM := linux
ifeq (Darwin, $(findstring Darwin, $(shell uname -a)))
  PLATFORM = darwin
endif

build:
	@CGO_ENABLED=0 go build -o "dist/app" cmd/server/main.go

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
