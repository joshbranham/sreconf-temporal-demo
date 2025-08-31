
@PHONY: start-temporal worker

start-temporal:
	@temporal server start-dev
worker:
	@go run main.go
