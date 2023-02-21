#!/bin/sh
go mod init github.com/jaaumeadrover/3aPracticaConcurrent
go get github.com/streadway/amqp
go get github.com/rabbitmq/amqp091-go
docker-compose up
