//AUTORS: DIEGO BERMEJO CABAÑAS I MARC CAÑELLAS GOMEZ
//VIDEO EXPLICATIU: https://www.youtube.com/watch?v=7js7pzTPZSY

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
	PROCESADORES = 6
	DIAL         = "amqp://guest:guest@localhost:5672/"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Funció que simula un client que menja una determinada quantitat de peces de sushi
func client(nsushis int, ch *amqp.Channel, q amqp.Queue, p amqp.Queue, done chan Empty) {
	//miram si a la cua de permisos hi ha 10 missatjes per a començar a menjar
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

	//bucle per a menjar les peces de sushi pasades per paràmetre
	for i := 0; i < nsushis; i++ {

		//esperam a tenir permisos
		mssg, _, err := ch.Get(p.Name, true)
		failOnError(err, "Failed to get a channel message")
		for mssg.Body == nil {
			mssg, _, err = ch.Get(p.Name, true)
			failOnError(err, "Failed to get a channel message")
		}

		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)

		//obtenim els missatjes de la cua de sushis per a menjar-los
		d, _, err1 := ch.Get(q.Name, true)
		failOnError(err1, "Failed to get a channel message")

		log.Printf("He menjat %s", d.Body)

		//miram la quantitat de peces que hi ha al plat un cop hem menjat una
		peces, err := ch.QueueInspect(q.Name)
		failOnError(err, "Failed to connect to RabbitMQ")

		log.Printf("Al plat hi ha %d peces", peces.Messages)

		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	}

	log.Printf("He acabat")

	done <- Empty{} //acaba el client
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

	//cua amb els permisos
	p, err := ch.QueueDeclare(
		"permisos", // name
		true,       // durablel
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

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

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	//calculam la quantitat de peces de suhsi que es menjarà
	rand.Seed(time.Now().UnixNano())
	nsushis := rand.Intn(15)

	go client(nsushis+1, ch, q, p, done)
	<-done

}
