package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rw)
		outcomeStr := ""

		gamelog := pubsub.GameLog{
			CurrentTime: time.Now(),
			Message:     outcomeStr,
			Username:    gs.Player.Username,
		}

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeYouWon:
			outcomeStr = fmt.Sprintf("%s won a war against %s", winner, loser)
			pubsub.PublishGameLog(ch, gamelog)
			return pubsub.Ack
		case gamelogic.WarOutcomeOpponentWon:
			outcomeStr = fmt.Sprintf("%s won a war against %s", winner, loser)
			pubsub.PublishGameLog(ch, gamelog)
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			outcomeStr = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			pubsub.PublishGameLog(ch, gamelog)
			return pubsub.Ack
		default:
			fmt.Printf("Unknown war outcome: %d\n", outcome)
			return pubsub.NackDiscard
		}
	}
}
