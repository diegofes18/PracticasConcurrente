package main

import (
	"log"
        "time"
        "fmt"
        "math/rand"
	amqp "github.com/streadway/amqp"
)
const (
        Productores = 4
        Consumidores = 4
        Nmensajes = 4

)

var prodnames = [Productores]string{"Ferran","Miquel","Tomeu","Toni"}
var consnames = [Consumidores]string{"Sara", "Maria","Gloria","Xisca"}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
type Empty struct{}

func producer(id int, name string,ch *amqp.Channel, err error, q amqp.Queue, done chan Empty){
        for i := 0; i < Nmensajes; i++ {
                body := fmt.Sprintf("element %d, from %s", i,name)
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
		log.Printf(" [*] Producer %d , %s sends %s", id, name, body)

		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
        }
        done <- Empty{}
}

func main() {
	
        conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable (the queue will survive a broker restart)
		false,   // delete when unused
		false,   // exclusive (used by only one connection and the queue will be deleted when that connection closes)
		false,   // no-wait (the server will not respond to the method. The client should not wait for a reply method)
		nil,     // arguments (Those are provided by clients when they declare queues (exchanges) and control various optional features, such as queue length limit or TTL.)
	)
	failOnError(err, "Failed to declare a queue")
        done := make(chan Empty, 1)
	for i := 0; i < Productores; i++ {
		go producer(i, prodnames[i], ch, err, q, done)
	}
	for i := 0; i < Productores; i++ {
		<-done
	}
}
