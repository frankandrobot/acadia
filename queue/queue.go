package queue

import (
	"github.com/frankandrobot/acadia/messaging"
)

func Queue() chan messaging.Payload {
	queue := make(chan messaging.Payload)
	go func() {
		for {
			payload := <-queue
			result := payload.Action()
			payload.Done <- result
		}
	}()
	return queue
}

type Contents struct {
	Contents string `form:"contents" json:"contents" binding:"required"`
}
