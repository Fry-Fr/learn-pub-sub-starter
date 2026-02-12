package main

import (
	"fmt"
	"os"
	"os/signal"

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("\nShutting down gracefully...")

}
