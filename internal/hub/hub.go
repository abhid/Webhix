package hub

import "sync"

// Hub is an in-memory pub/sub that fans out webhook events to SSE subscribers
// grouped by token.
type Hub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan []byte
	done        chan struct{}
}

// New returns an initialised Hub.
func New() *Hub {
	return &Hub{
		subscribers: make(map[string][]chan []byte),
		done:        make(chan struct{}),
	}
}

// Done returns a channel that is closed when the hub shuts down.
// SSE handlers should select on this to exit cleanly during graceful shutdown.
func (h *Hub) Done() <-chan struct{} {
	return h.done
}

// Subscribe registers a new subscriber for the given token.
// It returns a receive-only channel that will receive published payloads and an
// unsubscribe function that must be called when the subscriber is done (e.g.
// when the SSE client disconnects).
func (h *Hub) Subscribe(token string) (<-chan []byte, func()) {
	ch := make(chan []byte, 16)

	h.mu.Lock()
	h.subscribers[token] = append(h.subscribers[token], ch)
	h.mu.Unlock()

	unsubscribe := func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		subs := h.subscribers[token]
		for i, sub := range subs {
			if sub == ch {
				h.subscribers[token] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		if len(h.subscribers[token]) == 0 {
			delete(h.subscribers, token)
		}
		// ch is intentionally not closed here: Publish holds a copied reference
		// and a non-blocking send to a closed channel panics. The channel is
		// simply dropped from the map and will be GC'd when no longer referenced.
	}

	return ch, unsubscribe
}

// Close signals all SSE handlers to exit and clears the subscriber map.
// It does NOT close subscriber channels — that would race with Publish.
func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	select {
	case <-h.done:
		// already closed; nothing to do
	default:
		close(h.done)
	}

	for token := range h.subscribers {
		delete(h.subscribers, token)
	}
}

// Publish sends data to all current subscribers of the given token.
// Slow subscribers are skipped (non-blocking send) to avoid blocking the
// webhook handler.
func (h *Hub) Publish(token string, data []byte) {
	h.mu.RLock()
	subs := h.subscribers[token]
	// Copy the slice so we can release the lock before sending.
	targets := make([]chan []byte, len(subs))
	copy(targets, subs)
	h.mu.RUnlock()

	for _, ch := range targets {
		select {
		case ch <- data:
		default:
			// subscriber is too slow; drop the event for this subscriber
		}
	}
}
