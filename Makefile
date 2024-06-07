producer:
	@go run src/producer/producer.go

consumer:
	@go run src/consumer/consumer.go

message-queue:
	@docker run -d \
		-p 5672:5672 \
		-p 15672:15672 \
		-e RABBITMQ_DEFAULT_USER=rabbitmq \
		-e RABBITMQ_DEFAULT_PASS=rabbitmq \
		--name rabbitmq \
		rabbitmq:3.13-rc-management

minikube-delete-restart:
	@minikube delete
	@minikube start

fmt:
	@go fmt ./...

check:
	@golangci-lint run

generate:
	@go generate ./...

test: generate
	@go test ./... -cover -v -race

.PHONY: fmt check producer consumer message-queue minikube-delete-restart generate test