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

const webPort = "80"

type Config struct {
	RabbitConn *amqp.Connection
}

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := rabbitConnect()
	if err != nil {
		log.Println("Error connecting to RabbitMQ")
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		RabbitConn: rabbitConn,
	}

	log.Printf("Starting Broker service on port %s\n", webPort)

	// Define HTTP Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// Start the Server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func rabbitConnect() (*amqp.Connection, error) {
	var counts int64
	var connection *amqp.Connection

	// don't return until rabbit is ready
	for {
		username := os.Getenv("RABBITMQ_USERNAME")
		password := os.Getenv("RABBITMQ_PASSWORD")
		c, err := amqp.Dial("amqp://" + username + ":" + password + "@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			log.Println("Connected to RabbitMQ!")
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff := time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Backing off for %s seconds", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
