package pulse

import "sync/atomic"

// A pulse is meant to be a pub-sub signaler
// basically I want to get multiple reads from a single channel write
type Pulse struct {
	signal      chan bool
	subscribers []chan bool
	len         int32
}

// Init new pulse with signal channel
func NewPulse(signal chan bool) *Pulse {
	pulse := &Pulse{
		signal: signal,
		len:    0,
	}
	pulse.listen()
	return pulse
}

// subscribes to the signal bit
// if anything is written to the signal chan
// the returned chan will be written to
func (p *Pulse) Subscribe() chan bool {
	// add new chan to subscribers
	newChan := make(chan bool)
	p.subscribers = append(p.subscribers, newChan)
	atomic.AddInt32(&p.len, 1)
	return newChan
}

func (p *Pulse) listen() {
	go func() {
		for range p.signal {
			if p.len == 0 {
				continue
			}
			// this will ensure we read the chans in order
			subChan := make(chan chan bool)
			for _, sub := range p.subscribers {
				go func() {
					mySub := <-subChan
					mySub <- true
				}()
				subChan <- sub
			}
		}
	}()
}
