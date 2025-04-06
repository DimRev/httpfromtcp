.PHONY: run-tcplistener

run-tcplistener:
	@go run ./cmd/tcplistener

run-udpsender:
	@go run ./cmd/udpsender