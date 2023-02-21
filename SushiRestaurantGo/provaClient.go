package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Estructura vacia
type Empty struct{}

const (
	PROCESADORES = 4
	NPIEZAS      = 10
	DIAL         = "amqp://guest:guest@localhost:5672/"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func client(nsushis int, ch *amqp.Channel, q amqp.Queue, p amqp.Queue, done chan Empty) {
	permiso, err := ch.QueueInspect(p.Name)
	if err != nil {
		panic(err)
	}
	if permiso.Messages < 10 {
		_, err := ch.QueuePurge(p.Name, false)
		if err != nil {
			panic(err)
		}
	}

	log.Printf("Bon vespre, vinc a sopar de suhi")
	log.Printf("Avui menjaré %d peces", nsushis)

	for i := 0; i < nsushis; i++ {

		mssg, _, err := ch.Get("permisos", true)
		failOnError(err, "Failed to get a channel message")
		for mssg.Body == nil {
			mssg, _, err = ch.Get("permisos", true)
			failOnError(err, "Failed to get a channel message")
		}

		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)

		d, _, err1 := ch.Get(q.Name, true)
		failOnError(err1, "Failed to get a channel message")

		log.Printf("He menjat %s", d.Body)

		peces, err := ch.QueueInspect(q.Name)
		failOnError(err, "Failed to connect to RabbitMQ")

		log.Printf("Al plat hi ha %d peces", peces.Messages)

		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	}

	log.Printf("He acabat")
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

	p, err := ch.QueueDeclare(
		"permisos", // name
		true,       // durablel
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	q, err := ch.QueueDeclare(
		"sushis", // name
		true,     // durablel
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	rand.Seed(time.Now().UnixNano())
	nsushis := rand.Intn(15)

	go client(nsushis+1, ch, q, p, done)
	<-done
	// Salir cuando se presione CTRL+C

}
