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
	fmt.Println("Starting Peril client...")

	// connect to RabbitMQ
	const rabbitMqConnection = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitMqConnection)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMq: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error getting username: %v", err)
	}

	queueName := routing.PauseKey + "." + username
	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gameState := gamelogic.NewGameState(username)

	// main client loop
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			if err = gameState.CommandSpawn(words); err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			if _, err = gameState.CommandMove(words); err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gameState.CommandStatus()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command: " + words[0])
		}
	}
}
