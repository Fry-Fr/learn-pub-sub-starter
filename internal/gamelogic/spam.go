package gamelogic

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CommandSpam(gs GameState, ch *amqp.Channel, n int) error {
	for i := 0; i < n; i++ {
		msg := GetMaliciousLog()
		key := fmt.Sprintf("%s.%s", routing.GameLogSlug, gs.GetUsername())
		err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, key, routing.GameLog{
			CurrentTime: time.Now(),
			Message:     msg,
			Username:    gs.GetUsername(),
		})
		if err != nil {
			return fmt.Errorf("failed to publish spam message: %w", err)
		}
	}
	return nil
}
