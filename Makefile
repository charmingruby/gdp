.PHONY: run-client
run-client:
	go run ./cmd/client/main.go

.PHONY: run-server
run-server:
	go run ./cmd/server/main.go
