package chromecast

import (
	"context"
)

type RequestCounter struct {
	ctx         context.Context
	needCounter chan chan int
	id          int
}

func NewRequestCounter(ctx context.Context) *RequestCounter {
	return &RequestCounter{
		ctx:         ctx,
		needCounter: make(chan chan int, 100),
		id:          1,
	}
}

func (reqCounter *RequestCounter) GetRequestCounter() chan int {
	returnChan := make(chan int, 1)

	reqCounter.needCounter <- returnChan

	return returnChan
}

func (reqCounter *RequestCounter) Start() {
	go reqCounter.Listen()
}

func (reqCounter *RequestCounter) Listen() {
	for {
		select {
		case <-reqCounter.ctx.Done():
			return
		case sub := <-reqCounter.needCounter:
			reqCounter.id++
			sub <- reqCounter.id
		}
	}
}
