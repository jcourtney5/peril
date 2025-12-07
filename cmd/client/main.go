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

	gs := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	// main client loop
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			if err = gs.CommandSpawn(words); err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			if _, err = gs.CommandMove(words); err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gs.CommandStatus()
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
