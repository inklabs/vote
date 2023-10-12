package listener

import (
	"log"

	"github.com/inklabs/vote/event"
)

type ElectionWinnerVoterNotification struct{}

func (e *ElectionWinnerVoterNotification) On(event event.ElectionWinnerWasSelected) error {
	log.Printf("%+v", event)
	return nil
}
