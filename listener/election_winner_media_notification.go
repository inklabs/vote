package listener

import (
	"context"
	"time"

	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/pkg/sleep"
)

type ElectionWinnerMediaNotification struct{}

func NewElectionWinnerMediaNotification() *ElectionWinnerMediaNotification {
	return &ElectionWinnerMediaNotification{}
}

func (e *ElectionWinnerMediaNotification) On(ctx context.Context, _ event.ElectionWinnerWasSelected) error {
	_, span := tracer.Start(ctx, "vote.send-media-notification")
	defer span.End()

	sleep.Rand(1 * time.Millisecond)
	return nil
}
