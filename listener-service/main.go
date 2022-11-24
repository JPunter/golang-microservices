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
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println("Error connecting to RabbitMQ")
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
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
		log.Printf("Backing off for %s seconds", backoff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
