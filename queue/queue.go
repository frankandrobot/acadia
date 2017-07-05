package queue

import (
	"github.com/frankandrobot/acadia/messaging"
)

type Payload struct {
	Message messaging.Message
	Action  messaging.Action
	Done    chan messaging.ChanResult
}

type Queue chan Payload

func MakeQueue() Queue {
	queue := make(chan Payload)
	go func() {
		for {
			payload := <-queue
			result := payload.Action()
			payload.Done <- result
		}
	}()
	return queue
}

func (q Queue) Add(action messaging.Action) messaging.ChanResult {
	payload := Payload{
		Action: action,
		Done:   make(chan messaging.ChanResult),
	}
	go func() { q <- payload }()
	return <-payload.Done
}
