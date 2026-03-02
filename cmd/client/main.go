package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/cmd/client/handlers"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ampqURL := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(ampqURL)
	if err != nil {
		fmt.Printf("Failed to connect to Peril client: %s\n", err)
		return
	}
	defer connection.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("Failed to welcome client: %s\n", err)
		return
	}
	exchange := routing.ExchangePerilDirect
	routingKey := routing.PauseKey
	queueName := fmt.Sprintf("%s.%s", routingKey, username)
	gameState := gamelogic.NewGameState(username)

	// subscribe to pause/resume messages for this client
	err = pubsub.SubscribeJSON(connection, exchange, queueName, routingKey, pubsub.TransientQueue, handlers.HandlerPause(gameState))
	if err != nil {
		fmt.Printf("Failed to subscribe to pause messages: %s\n", err)
		return
	}

	for {
		inputWords := gamelogic.GetInput()
		if len(inputWords) == 0 {
			continue
		}
		command := inputWords[0]

		switch command {
		case "spawn":
			err = gameState.CommandSpawn(inputWords)
			if err != nil {
				fmt.Printf("Error processing command: %s\n", err)
			}
		case "move":
			_, err = gameState.CommandMove(inputWords)
			if err != nil {
				fmt.Printf("Error processing command: %s\n", err)
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Printf("Unknown command: %s\n", command)
			continue
		}
	}
}
