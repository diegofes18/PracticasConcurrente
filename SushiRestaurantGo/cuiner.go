//AUTORS: DIEGO BERMEJO CABAÑAS I MARC CAÑELLAS GOMEZ
//VIDEO EXPLICATIU: https://www.youtube.com/watch?v=7js7pzTPZSY

package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Estructura buida
type Empty struct{}

// Estructura que conté el tipus de sushi i la quantitat de peces de cada un
type PiezaSushi struct {
	tipo string
	n    int
}

const (
	DIAL = "amqp://guest:guest@localhost:5672/"
)

// Tipus de suhi
var sushis = []string{"Nigiri salmon", "Sashimi Tonyina", "Maki de cranc"}

// plat que conté els sushis
var plat = []PiezaSushi{}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Funció que s'encarrega de posar en el plat 10 peces de sushi aleatories
func cuiner(ch *amqp.Channel, q amqp.Queue, p amqp.Queue, done chan Empty) {
	//Buidam les cues de permisos i sushis per a tornar a posar 10 peces de suhi
	// i 10 permisos per a menjar-les
	_, err := ch.QueuePurge(p.Name, false)
	if err != nil {
		panic(err)
	}

	_, err = ch.QueuePurge(q.Name, false)
	if err != nil {
		panic(err)
	}

	fmt.Printf("El cuiner de sushi ja és aquí\n")
	fmt.Printf("El cuiner prepararà un plat amb: \n")

	//Random seed
	rand.Seed(time.Now().UTC().UnixNano())

	total := 0 //nombre de peces que hi ha al plat

	//Bucle per emplenar el plat amb els dos primers tipus de sushi
	for i := 0; i < len(sushis)-1; i++ {
		nPiezas := rand.Intn(10 - total)
		total = total + nPiezas
		fmt.Printf("%d peces de %s\n", nPiezas, sushis[i])
		plat = append(plat, PiezaSushi{sushis[i], nPiezas}) //afegim les peces al plat
	}

	//Calculam la quantitat de peces del tercer tipus de sushi
	i := 10 - total
	fmt.Printf("%d peces de %s\n", i, sushis[2])
	plat = append(plat, PiezaSushi{sushis[2], i}) //les afegim

	//bucle per afegir les peces de sushi a la cua
	for i := 0; i < len(sushis); i++ { //afegim totes les peces de cada tipus juntes
		for j := 0; j < plat[i].n; j++ {
			body := fmt.Sprintf("%s", plat[i].tipo)
			err := ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})

			failOnError(err, "Failed to publish a message")

			tiempoAleatorio := rand.Intn(2000)
			time.Sleep(time.Duration(tiempoAleatorio) * time.Millisecond)

			log.Printf(" [x] Posa dins el plat %s", body)
		}
	}

	fmt.Printf("Podeu menjar")

	//bucle per afegir a la cua de permisos els 10 permisos per a cada sushi
	for i := 0; i < len(sushis); i++ {
		for j := 0; j < plat[i].n; j++ {

			err := ch.Publish(
				"",     // exchange
				p.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(""),
				})
			failOnError(err, "Failed to publish a message")

		}
	}

	done <- Empty{} //el cuiner acaba

}

func main() {
	runtime.GOMAXPROCS(1)

	done := make(chan Empty, 1)

	conn, err := amqp.Dial(DIAL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//cua amb el plat de sushis
	q, err := ch.QueueDeclare(
		"sushis", // name
		true,     // durablel
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//cua amb els permisos per a menjar les peces de sushi
	p, err := ch.QueueDeclare(
		"permisos", // name
		true,       // durablel
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)

	failOnError(err, "Failed to declare a queue")

	go cuiner(ch, q, p, done)

	<-done

}
