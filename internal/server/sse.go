// Package server provides HTTP file server and dynamic rendering functionality.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// SSEBroadcaster manages Server-Sent Events connections for live reload
type SSEBroadcaster struct {
	clients   map[chan string]struct{}
	mu        sync.RWMutex
	messageCh chan string
	stopCh    chan struct{}
}

// NewSSEBroadcaster creates a new SSE broadcaster
func NewSSEBroadcaster() *SSEBroadcaster {
	return &SSEBroadcaster{
		clients:   make(map[chan string]struct{}),
		messageCh: make(chan string, 10),
		stopCh:    make(chan struct{}),
	}
}

// Start begins the broadcaster goroutine
func (b *SSEBroadcaster) Start() {
	go b.broadcastLoop()
}

// Stop stops the broadcaster
func (b *SSEBroadcaster) Stop() {
	close(b.stopCh)
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.clients {
		close(ch)
	}
	b.clients = make(map[chan string]struct{})
}

// Broadcast sends a message to all connected clients
func (b *SSEBroadcaster) Broadcast(eventType string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, string(jsonData))

	select {
	case b.messageCh <- message:
	default:
		// Channel full, drop message
	}
}

// broadcastLoop sends messages to all connected clients
func (b *SSEBroadcaster) broadcastLoop() {
	for {
		select {
		case <-b.stopCh:
			return
		case msg := <-b.messageCh:
			b.mu.RLock()
			for ch := range b.clients {
				select {
				case ch <- msg:
				default:
					// Client is slow, skip
				}
			}
			b.mu.RUnlock()
		}
	}
}

// addClient registers a new client
func (b *SSEBroadcaster) addClient(ch chan string) {
	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()
}

// removeClient unregisters a client
func (b *SSEBroadcaster) removeClient(ch chan string) {
	b.mu.Lock()
	delete(b.clients, ch)
	b.mu.Unlock()
}

// Handler returns an HTTP handler for SSE connections
func (b *SSEBroadcaster) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Create a channel for this client
		clientCh := make(chan string, 10)
		b.addClient(clientCh)
		defer b.removeClient(clientCh)

		// Get flusher
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		// Send initial connection message
		_, _ = fmt.Fprintf(w, "event: connected\ndata: {}\n\n")
		flusher.Flush()

		// Listen for messages or client disconnect
		for {
			select {
			case <-r.Context().Done():
				return
			case msg, ok := <-clientCh:
				if !ok {
					return
				}
				_, err := fmt.Fprint(w, msg)
				if err != nil {
					return
				}
				flusher.Flush()
			}
		}
	}
}
