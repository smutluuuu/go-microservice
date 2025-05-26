package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	rabbitCon, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitCon.Close()

	app := &Config{
		Rabbit: rabbitCon,
	}
	log.Printf("Starting broker service on port %s\n", webPort)

	// HTTP SERVER DEFINITION
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	// START THE SERVER
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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
