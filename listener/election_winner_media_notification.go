package listener

import (
	"time"

	"github.com/inklabs/vote/event"
)

type ElectionWinnerMediaNotification struct{}

func NewElectionWinnerMediaNotification() *ElectionWinnerMediaNotification {
	return &ElectionWinnerMediaNotification{}
}

func (e *ElectionWinnerMediaNotification) On(_ event.ElectionWinnerWasSelected) error {
	time.Sleep(2 * time.Millisecond)
	return nil
}
