SHELL := /bin/bash

gen: gen_sql gen_mocks

gen_sql:
	sqlc generate -f ./postgres-driver/sqlc/sqlc.yaml

gen_mocks:
	mockery --name=Driver --recursive --inpkg --case=underscore
	mockery --name=RelayWriter --recursive --inpkg --case=underscore
	mockery --name=ServiceRecordWriter --recursive --inpkg --case=underscore

test: test_unit test_env_up run_driver_tests test_env_down

test_unit:
	@echo "🛠 Running Unit tests..."
	@go test ./... -short || true
	@echo "✅ Unit tests completed!"

test_env_up:
	@echo "🧪 Starting up Transaction DB test database ..."
	@docker-compose -f ./testdata/docker-compose.test.yml up -d --remove-orphans --build
	@echo "⏳ Waiting for test DB to be ready ..."
	@attempts=0; while ! pg_isready -h localhost -p 5432 -U postgres -d postgres >/dev/null && [[ $$attempts -lt 5 ]]; do sleep 1; attempts=$$(($$attempts + 1)); done
	@[[ $$attempts -lt 5 ]] && echo "🐘 Test Transation DB is up ..." || (echo "❌ Test Transation DB failed to start" && make test_env_down >/dev/null && exit 1)
	@echo "🚀 Test environment is up ..."
test_env_down:
	@echo "🧪 Shutting down Portal HTTP DB test environment ..."
	@docker-compose -f ./testdata/docker-compose.test.yml down --remove-orphans >/dev/null
	@echo "✅ Test environment is down."

run_driver_tests:
	@echo "🚗 Running PGDriver tests..."
	@go test ./... -run Test_RunPGDriverSuite -count=1 || true
	@echo "✅ PGDriver tests completed!"
run_driver_tests_ci:
	go test ./... -run Test_RunPGDriverSuite -count=1

init-pre-commit:
	wget https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz
	python3 pre-commit-2.20.0.pyz install
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
	python3 pre-commit-2.20.0.pyz run --all-files
	rm pre-commit-2.20.0.pyz
