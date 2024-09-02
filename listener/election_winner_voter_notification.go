package listener

import (
	"time"

	"github.com/inklabs/vote/event"
)

type ElectionWinnerVoterNotification struct{}

func NewElectionWinnerVoterNotification() *ElectionWinnerVoterNotification {
	return &ElectionWinnerVoterNotification{}
}

func (e *ElectionWinnerVoterNotification) On(_ event.ElectionWinnerWasSelected) error {
	time.Sleep(2 * time.Millisecond)
	return nil
}
