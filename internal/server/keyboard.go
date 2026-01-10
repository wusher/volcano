// Package server provides HTTP file server and dynamic rendering functionality.
package server

import (
	"fmt"
	"io"
	"os"
	"sync"

	"golang.org/x/term"
)

// KeyboardHandler handles keyboard input for the dynamic server
type KeyboardHandler struct {
	writer       io.Writer
	stopCh       chan struct{}
	stdinFd      int
	oldState     *term.State
	mu           sync.Mutex
	running      bool
	onKeyPressed func(key rune)
}

// NewKeyboardHandler creates a new keyboard handler
func NewKeyboardHandler(writer io.Writer, onKeyPressed func(key rune)) *KeyboardHandler {
	return &KeyboardHandler{
		writer:       writer,
		stopCh:       make(chan struct{}),
		stdinFd:      int(os.Stdin.Fd()),
		onKeyPressed: onKeyPressed,
	}
}

// Start begins listening for keyboard input in a goroutine
func (k *KeyboardHandler) Start() error {
	// Check if stdin is a terminal
	if !term.IsTerminal(k.stdinFd) {
		return nil // Not a terminal, skip keyboard handling
	}

	// Put terminal into raw mode
	oldState, err := term.MakeRaw(k.stdinFd)
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}

	k.mu.Lock()
	k.oldState = oldState
	k.running = true
	k.mu.Unlock()

	go k.readLoop()
	return nil
}

// Stop stops the keyboard handler and restores terminal state
func (k *KeyboardHandler) Stop() {
	k.mu.Lock()
	defer k.mu.Unlock()

	if !k.running {
		return
	}

	k.running = false
	close(k.stopCh)

	// Restore terminal state
	if k.oldState != nil {
		_ = term.Restore(k.stdinFd, k.oldState)
	}
}

// readLoop continuously reads keyboard input
func (k *KeyboardHandler) readLoop() {
	buf := make([]byte, 1)
	for {
		select {
		case <-k.stopCh:
			return
		default:
			n, err := os.Stdin.Read(buf)
			if err != nil {
				return
			}
			if n > 0 {
				key := rune(buf[0])
				// Handle Ctrl+C (ETX, 0x03)
				if key == 3 {
					k.Stop()
					return
				}
				if k.onKeyPressed != nil {
					k.onKeyPressed(key)
				}
			}
		}
	}
}
