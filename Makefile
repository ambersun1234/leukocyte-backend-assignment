PRODUCER_IMAGE=assignment-producer
CONSUMER_IMAGE=assignment-consumer

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
	@minikube start --memory 8192	
	@minikube image load $(PRODUCER_IMAGE)
	@minikube image load $(CONSUMER_IMAGE)
	@kubectl apply -f ./kubernetes/configmaps
	@kubectl apply -f ./kubernetes/secrets
	@kubectl apply -f ./kubernetes/permissions
	@kubectl apply -f ./kubernetes/deployments
	@kubectl apply -f ./kubernetes/services
	
.PHONY: 
	fmt 
	check 
	generate 
	test 
	build-image 
	deploy-minikube