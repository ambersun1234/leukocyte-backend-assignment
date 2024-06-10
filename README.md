# Leukocyte Backend Assignment
[![CI test](https://github.com/ambersun1234/leukocyte-backend-assignment/actions/workflows/ci.yaml/badge.svg)](https://github.com/ambersun1234/leukocyte-backend-assignment/actions/workflows/ci.yaml)

## Introduction
This repository implements a simple distributed producer consumer pattern.\
The producer will produce a "task" to message queue every 10 seconds, and is consumed by the consumer\
The consumer will schedule the task to orchestration tool like k8s

## Spec
Please refer to [spec](./面試作業.pdf) for more detail

## Prerequisites
+ [Docker](https://www.docker.com/)
+ [Kubectl](https://kubernetes.io/docs/reference/kubectl/)
+ [Minikube](https://minikube.sigs.k8s.io/docs/start/?arch=%2Flinux%2Fx86-64%2Fstable%2Fbinary+download)

## Config
You can change the config in [./config.yaml](./config.yaml), or keep the default one

## Run
### Manual
```shell
$ make message-queue
$ make minikube-delete-restart
$ make producer
$ make consumer
```

> This approach will run the application code outside kubernetes

### Auto
```shell
$ make deploy-minikube
```

## Test
```shell
$ make test
```

## Check List
+ [x] producer
    + [x] should publish message to message queue
+ [x] consumer
    + [x] should acknowledge message
    + [x] should consume message from message queue
    + [x] should re-enqueue failed message
    + [x] should schedule job to kubernetes
+ [x] message queue
    + [x] should implement auto-recovery mechanism
+ [x] kubernetes
    + [x] should schedule job from job queue
+ [x] unit test
+ [x] minikube

## Author
+ [ambersun1234](https://github.com/ambersun1234)

## License
This project is licensed under MIT. - see the [LICENSE](./LICENSE) file for more detail
