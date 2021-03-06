package chanutil

import (
	"errors"
	"fmt"
	"time"
)

// Non-blocking channel.
type NBChan struct {
	ch        chan interface{}
	LogString string
}

func NewNBChan() *NBChan {
	return NewNBChan2(0, "nbchan")
}
func NewNBChan2(n int, logS string) *NBChan {
	ch := &NBChan{
		ch:        make(chan interface{}, n),
		LogString: logS,
	}
	return ch
}

//----------

// Non-blocking send now, or fails with error.
func (ch *NBChan) Send(v interface{}) error {
	select {
	case ch.ch <- v:
		return nil
	default:
		return errors.New("failed to send")
	}
}

// Receives or fails after timeout.
func (ch *NBChan) Receive(timeout time.Duration) (interface{}, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil, fmt.Errorf("%v: receive timeout", ch.LogString)
	case v := <-ch.ch:
		return v, nil
	}
}

//----------

func (ch *NBChan) NewBufChan(n int) {
	ch.ch = make(chan interface{}, n)
}

// Setting the channel to zero allows a send to fail immediatly if there is no receiver waiting.
func (ch *NBChan) SetBufChanToZero() {
	ch.NewBufChan(0)
}
