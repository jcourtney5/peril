package main

import (
	"fmt"
	"log"

	"github.com/jcourtney5/peril/internal/pubsub"
	"github.com/jcourtney5/peril/internal/routing"
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

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening channel: %v", err)
	}

	err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		log.Fatalf("Error publishing JSON: %v", err)
	}

	fmt.Println("Pause message sent!")
}
