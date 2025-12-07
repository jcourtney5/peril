package main

import (
	"fmt"
	"log"

	"github.com/jcourtney5/peril/internal/gamelogic"
	"github.com/jcourtney5/peril/internal/pubsub"
	"github.com/jcourtney5/peril/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	// connect to RabbitMQ
	const rabbitMqConnection = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitMqConnection)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMq: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	// declare and bind the peril_topic queue
	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatalf("could not declare %v queue: %v", routing.GameLogSlug, err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening channel: %v", err)
	}

	gamelogic.PrintServerHelp()

	// Main loop
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			sendPause(publishCh)
		case "resume":
			sendResume(publishCh)
		case "quit":
			fmt.Println("Quitting the game...")
			return
		default:
			fmt.Println("Unknown command: " + words[0])
		}

	}
}

func sendPause(publishCh *amqp.Channel) {
	fmt.Println("Sending pause message")

	err := pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		log.Fatalf("Error sending pause message: %v", err)
	}

	fmt.Println("Pause message sent!")
}

func sendResume(publishCh *amqp.Channel) {
	fmt.Println("Sending resume message")

	err := pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: false,
	})
	if err != nil {
		log.Fatalf("Error sending resume message: %v", err)
	}

	fmt.Println("Resume message sent!")
}
