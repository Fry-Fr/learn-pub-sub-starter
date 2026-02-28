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
		fmt.Printf("Failed to connect to RabbitMQ: %s\n", err)
		return
	}
	defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ")
	fmt.Println("Starting Peril server...")

	// make sure we have a durable queue for game logs attached to the topic exchange
	exchange := routing.ExchangePerilTopic
	queueName := routing.GameLogSlug
	routingKey := fmt.Sprintf("%s.*", routing.GameLogSlug)
	// durable queues survive restarts, are not auto‑deleted, and allow multiple consumers
	_, _, err = pubsub.DeclareAndBind(connection, exchange, queueName, routingKey, pubsub.DurableQueue)
	if err != nil {
		fmt.Printf("Failed to declare and bind game_logs queue: %s\n", err)
		return
	}

	gamelogic.PrintServerHelp()
	for {
		inputSlice := gamelogic.GetInput()
		if len(inputSlice) == 0 {
			continue
		}
		command := inputSlice[0]
		switch command {
		case "pause":
			fmt.Println("sending pause message...")
			channel, err := connection.Channel()
			if err != nil {
				fmt.Printf("Failed to open a channel: %s\n", err)
				return
			}
			defer channel.Close()

			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				fmt.Printf("Failed to publish pause message: %s\n", err)
				return
			}

		case "resume":
			fmt.Println("sending resume message...")
			channel, err := connection.Channel()
			if err != nil {
				fmt.Printf("Failed to open a channel: %s\n", err)
				return
			}
			defer channel.Close()

			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
			if err != nil {
				fmt.Printf("Failed to publish resume message: %s\n", err)
				return
			}

		case "quit":
			fmt.Println("Quitting the server...")
			return
		default:
			fmt.Printf("Unknown command: %s\n", command)
			continue
		}
	}
}
