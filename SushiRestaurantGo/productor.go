package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PRODUCTORES  = 2
	CONSUMIDORES = 2
	PROCESADORES = 4
	TO_PRODUCE   = 10
	DIAL         = "amqp://guest:guest@localhost:5672/"
	NOMBRE_COLA  = "buffer"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Empty struct{}

func productor(id int, ch *amqp.Channel, err error, q amqp.Queue, done chan Empty) {
	for i := 0; i < TO_PRODUCE; i++ {
		body := fmt.Sprintf("Productor %d produce el elmento %d", id, i)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Productor %d envia %s", id, body)
		tiempoAleatorio := rand.Intn(2000)
		log.Printf("Productor %d duerme %d ms", id, tiempoAleatorio)
		time.Sleep(time.Duration(tiempoAleatorio) * time.Millisecond)
	}
	done <- Empty{}
}

func main() {
	runtime.GOMAXPROCS(PROCESADORES)

	done := make(chan Empty, 1)

	conn, err := amqp.Dial(DIAL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		NOMBRE_COLA, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for i := 0; i < PRODUCTORES; i++ {
		go productor(i, ch, err, q, done)
	}

	for i := 0; i < PRODUCTORES; i++ {
		<-done
	}

}
