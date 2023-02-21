package main

import (
	"log"
	"runtime"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Empty struct{}

const (
	PROCESADORES = 4
	NPIEZAS      = 10
	DIAL         = "amqp://guest:guest@localhost:5672/"
	NOMBRE_COLA  = "buffer"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func gangster(ch *amqp.Channel, q amqp.Queue, done chan Empty) {

	queueInfo, err := ch.QueueInspect(q.Name)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bon dia, soc maleducat")
	log.Printf("Vull menjar tot el sushi del plat")

	// Consume messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	log.Printf("He menjat totes les peces, total %d", queueInfo.Messages)
	failOnError(err, "Failed to register a consumer")

	for i := 0; i < queueInfo.Messages; i++ {
		d := <-msgs
		d.Ack(false)
	}
	log.Printf("Romp el plat")
	done <- Empty{}
}

func main() {
	runtime.GOMAXPROCS(PROCESADORES)

	conn, err := amqp.Dial(DIAL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	done := make(chan Empty, 1)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		NOMBRE_COLA, // name
		true,        // durablel
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	go gangster(ch, q, done)
	<-done
}
