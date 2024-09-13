package listener

import (
	"context"
	"time"

	"github.com/inklabs/vote/event"
	"github.com/inklabs/vote/pkg/sleep"
)

type ElectionWinnerVoterNotification struct{}

func NewElectionWinnerVoterNotification() *ElectionWinnerVoterNotification {
	return &ElectionWinnerVoterNotification{}
}

func (e *ElectionWinnerVoterNotification) On(ctx context.Context, _ event.ElectionWinnerWasSelected) error {
	_, span := tracer.Start(ctx, "vote.send-voter-notification")
	defer span.End()

	sleep.Rand(1 * time.Millisecond)
	return nil
}
