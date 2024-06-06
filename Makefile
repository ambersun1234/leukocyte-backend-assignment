producer:
	go run src/producer/producer.go

consumer:
	go run src/consumer/consumer.go

message-queue:
	docker run -d \
		-p 5672:5672 \
		-p 15672:15672 \
		-e RABBITMQ_DEFAULT_USER=rabbitmq \
		-e RABBITMQ_DEFAULT_PASS=rabbitmq \
		--name rabbitmq \
		rabbitmq:3.13-rc-management

fmt:
	@go fmt ./...

check:
	@golangci-lint run