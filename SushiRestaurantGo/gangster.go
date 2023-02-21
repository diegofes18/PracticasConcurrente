//AUTORS: DIEGO BERMEJO CABAÑAS I MARC CAÑELLAS GOMEZ
//VIDEO EXPLICATIU: https://www.youtube.com/watch?v=7js7pzTPZSY

package main

import (
	"log"
	"runtime"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Empty struct{}

const (
	PROCESADORES = 4
	DIAL         = "amqp://guest:guest@localhost:5672/"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Funció que simula un client que buida la cua de sushis
func gangster(ch *amqp.Channel, q amqp.Queue, p amqp.Queue, done chan Empty) {
	// Obtenim la informació de la cua per a saber si hi ha 10 peces de suhi
	queue, err := ch.QueueInspect(q.Name)
	if err != nil {
		panic(err)
	}

	//buidam les cues ja que no pot començar a menjar
	if queue.Messages < 10 {
		_, err := ch.QueuePurge(p.Name, false)
		if err != nil {
			panic(err)
		}

		_, err = ch.QueuePurge(q.Name, false)
		if err != nil {
			panic(err)
		}

	}

	log.Printf("Bon vespre, vinc a sopar de sushi")
	log.Printf("Ho vull tot!")

	//Consume messages
	permiso_msg, err := ch.Consume(
		p.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//Consume messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//esperam a que el canal doni permís
	<-permiso_msg

	// Obtenim la informació de la cua per a saber la quantitat de peces que menja
	queue, err = ch.QueueInspect(q.Name)
	if err != nil {
		panic(err)
	}
	mensajes := queue.Messages + 1
	counter := 0
	//bucle on es consumeixen totes les peces de sushi
	for i := 0; i < mensajes; {
		d := <-msgs
		queue, err := ch.QueueInspect("sushis")
		if err != nil {
			panic(err)
		}
		mensajes = queue.Messages
		//log.Println(queue.Messages)
		d.Ack(false)
		counter = counter + 1

	}

	log.Printf("Agafa totes les peces, en total %d", counter)
	log.Printf("Romp el plat")
	log.Printf("Men vaig")

	done <- Empty{} //acaba el gangster
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

	//cua amb el plat de sushi
	q, err := ch.QueueDeclare(
		"sushis", // name
		true,     // durablel
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//cua amb els permisos per a començar a menjar
	p, err := ch.QueueDeclare(
		"permisos", // name
		true,       // durablel
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	go gangster(ch, q, p, done)
	<-done
}
