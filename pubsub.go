// Pub and sub using channels
package main

import (
	"sync"
)

type PubSub struct {
	subscribers map[string][]chan Match
	mu          sync.RWMutex
}

var ps = &PubSub{
	subscribers: make(map[string][]chan Match),
}

func (p *PubSub) Subscribe(matchID string) chan Match {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch := make(chan Match, 1)
	p.subscribers[matchID] = append(p.subscribers[matchID], ch)
	return ch
}

func (p *PubSub) Unsubscribe(matchId string, ch chan Match) {
	p.mu.Lock()
	defer p.mu.Unlock()

	subs := p.subscribers[matchId]

	for i, sub := range subs {
		if sub == ch {
			p.subscribers[matchId] = append(subs[:i], subs[i+1:]...)
			close(ch)
			return
		}
	}
}

func (p *PubSub) Publish(matchID string, match Match) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, ch := range p.subscribers[matchID] {
		select {
		case ch <- match:
		default:
			// skip if subscriber is slow
		}
	}
}
