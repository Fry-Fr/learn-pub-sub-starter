package pubsub

import (
	"bytes"
	"encoding/gob"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	// ensure the queue exists and is bound; the returned channel is closed by the helper
	_, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	// open a fresh channel for consuming
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	// note: we intentionally do not close the channel here because the consumer
	// goroutine will use it for the lifetime of the program. The caller is
	// responsible for closing the connection when the application exits.

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer (empty = auto-generated)
		false,     // autoAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // args
	)
	if err != nil {
		ch.Close()
		return err
	}

	// start goroutine to process deliveries
	go func() {
		for d := range msgs {
			var v T
			if err := gob.NewDecoder(bytes.NewReader(d.Body)).Decode(&v); err != nil {
				// ignore malformed messages
				continue
			}
			outcome := handler(v)
			switch outcome {
			case Ack:
				d.Ack(false)
			case NackRequeue:
				d.Nack(false, true)
			case NackDiscard:
				d.Nack(false, false)
			}
		}
	}()

	return nil
}
