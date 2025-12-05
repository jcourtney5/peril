package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to RabbitMQ
	const rabbitMqConnection = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitMqConnection)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMq: %v", err)
	}
	defer conn.Close()

	fmt.Println("Peril game server started and connected to RabbitMQ")

	// wait for Ctrl+C to exit program
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Server running, press Ctrl+C to stop...")
	<-signalChan // Will block here until user hits ctrl+c
	fmt.Println("Stopping server and disconnecting from RabbitMQ...")
}
