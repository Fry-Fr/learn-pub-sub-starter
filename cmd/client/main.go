package main

import (
	"fmt"
	"os"
	"os/signal"

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
	_, _, err = pubsub.DeclareAndBind(connection, exchange, queueName, routingKey, pubsub.TransientQueue)
	if err != nil {
		fmt.Printf("Failed to declare and bind queue: %s\n", err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("\nShutting down gracefully...")
}
