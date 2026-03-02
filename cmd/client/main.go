package main

import (
	"fmt"

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
	routingKey := routing.PauseKey
	queueName := fmt.Sprintf("%s.%s", routingKey, username)
	gameState := gamelogic.NewGameState(username)

	// subscribe to pause/resume messages for this client
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, queueName, routingKey, pubsub.TransientQueue, handlerPause(gameState))
	if err != nil {
		fmt.Printf("Failed to subscribe to pause messages: %s\n", err)
		return
	}

	armyMovesRoutingKey := fmt.Sprintf("%s.*", routing.ArmyMovesPrefix)
	armyMovesQueueName := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username)

	// subscribe to move messages for this client
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, armyMovesQueueName, armyMovesRoutingKey, pubsub.TransientQueue, handlerMove(gameState))
	if err != nil {
		fmt.Printf("Failed to subscribe to move messages: %s\n", err)
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
			channel, err := connection.Channel()
			if err != nil {
				fmt.Printf("Failed to open a channel: %s\n", err)
				return
			}
			defer channel.Close()
			armyMove, err := gameState.CommandMove(inputWords)
			if err != nil {
				fmt.Printf("Error processing command: %s\n", err)
			}
			err = pubsub.PublishJSON(channel, routing.ExchangePerilTopic, armyMovesRoutingKey, armyMove)
			if err != nil {
				fmt.Printf("Failed to publish move: %s\n", err)
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
