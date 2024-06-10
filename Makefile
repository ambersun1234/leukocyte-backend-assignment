PRODUCER_IMAGE=assignment-producer
CONSUMER_IMAGE=assignment-consumer

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

build-image: Dockerfile.producer Dockerfile.consumer
	@docker build -t $(PRODUCER_IMAGE) -f ./Dockerfile.producer .
	@docker build -t $(CONSUMER_IMAGE) -f ./Dockerfile.consumer .

deploy-minikube: build-image
	@minikube start
	@minikube image load $(PRODUCER_IMAGE)
	@minikube image load $(CONSUMER_IMAGE)
	@kubectl apply -f ./kubernetes/services
	@kubectl apply -f ./kubernetes/deployments
	@kubectl apply -f ./kubernetes/permissions
	@kubectl apply -f ./kubernetes/secrets

.PHONY: fmt check producer consumer message-queue minikube-delete-restart generate test build-image deploy-minikube