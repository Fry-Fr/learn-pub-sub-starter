package main

import (
	"fmt"
	"os"
	"os/signal"

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("\nShutting down gracefully...")

}
