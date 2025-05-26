package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	// try to connect to a rabbitmq server
	rabbitCon, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitCon.Close()

	// start listening for messages
	consumer, err := event.NewConsumer(rabbitCon)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var count int64
	var delay = 1 * time.Second
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			count++

		} else {
			fmt.Println("Connected to RabbitMQ")
			connection = conn
			break
		}
		if count > 5 {
			fmt.Println(err)
			return nil, err
		}
		delay = time.Duration(math.Pow(float64(count), 2)) * time.Second
		time.Sleep(delay)
		continue
	}

	return connection, nil
}
