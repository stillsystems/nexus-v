package cli

import (
	"fmt"
	"time"
)

type Spinner struct {
	stop   chan struct{}
	active bool
}

func NewSpinner() *Spinner {
	return &Spinner{stop: make(chan struct{})}
}

func (s *Spinner) Start(msg string) {
	s.active = true
	fmt.Print("\033[?25l") // Hide cursor
	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.stop:
				return
			default:
				fmt.Printf("\r%s %s", frames[i], msg)
				i = (i + 1) % len(frames)
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.active = false
	close(s.stop)
	fmt.Print("\r\033[K")   // Clear line
	fmt.Print("\033[?25h") // Show cursor
}

