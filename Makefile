gen_mocks:
	mockery --name=Driver --recursive --inpkg --case=underscore
	mockery --name=RelayWriter --recursive --inpkg --case=underscore
	mockery --name=ServiceRecordWriter --recursive --inpkg --case=underscore

init-pre-commit:
	wget https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz
	python3 pre-commit-2.20.0.pyz install
	python3 pre-commit-2.20.0.pyz autoupdate
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
	python3 pre-commit-2.20.0.pyz run --all-files
	rm pre-commit-2.20.0.pyz
