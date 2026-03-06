package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var gobBytes bytes.Buffer
	err := gob.NewEncoder(&gobBytes).Encode(val)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body:        gobBytes.Bytes(),
	})
	if err != nil {
		return err
	}
	return nil
}

type GameLog struct {
	CurrentTime time.Time
	Message     string
	Username    string
}

func PublishGameLog(ch *amqp.Channel, gl GameLog) error {
	routingKey := fmt.Sprintf("%s.%s", routing.GameLogSlug, gl.Username)
	exchange := routing.ExchangePerilTopic
	err := PublishGob(ch, exchange, routingKey, gl)
	if err != nil {
		return err
	}
	return nil
}
