package listener

import (
	"context"
	"time"

	"github.com/inklabs/vote/event"
)

type ElectionWinnerMediaNotification struct{}

func NewElectionWinnerMediaNotification() *ElectionWinnerMediaNotification {
	return &ElectionWinnerMediaNotification{}
}

func (e *ElectionWinnerMediaNotification) On(ctx context.Context, _ event.ElectionWinnerWasSelected) error {
	_, span := tracer.Start(ctx, "send-media-notification")
	defer span.End()

	time.Sleep(1 * time.Millisecond)
	return nil
}
