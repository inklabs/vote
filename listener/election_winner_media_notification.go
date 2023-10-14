package listener

import (
	"log"

	"github.com/inklabs/vote/event"
)

type ElectionWinnerMediaNotification struct{}

func NewElectionWinnerMediaNotification() *ElectionWinnerMediaNotification {
	return &ElectionWinnerMediaNotification{}
}

func (e *ElectionWinnerMediaNotification) On(event event.ElectionWinnerWasSelected) error {
	log.Printf("%#v", event)
	return nil
}
