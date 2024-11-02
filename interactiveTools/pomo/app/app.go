package app

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"haonguyen.tech/interactiveTools/pomo/pomodoro"
)

type App struct {
	ctx        context.Context
	controller *termdash.Controller
	terminal   *tcell.Terminal
	errorCh    chan error
	redrawCh   chan bool
	size       image.Point
}

func New(config *pomodoro.IntervalConfig) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	redrawCh := make(chan bool)
	errorCh := make(chan error)

	quitter := func(k *terminalapi.Keyboard) { // Quit on pressing 'q'
		if k.Key == 'q' || k.Key == 'Q' {
			cancel() // Cancels the context, exiting the app
		}
	}
	// Create the terminal.
	term, err := tcell.New()
	if err != nil {
		return nil, err
	}

	c, err := newGrid(ctx, term, config, errorCh, redrawCh)
	if err != nil {
		return nil, err
	}
	controller, err := termdash.NewController(term, c, termdash.KeyboardSubscriber(quitter))
	if err != nil {
		return nil, err
	}

	return &App{
		ctx:        ctx,
		controller: controller,
		terminal:   term,
		errorCh:    errorCh,
		redrawCh:   redrawCh,
	}, nil
}

func (a *App) resize() error {
	if a.size.Eq(a.terminal.Size()) {
		return nil
	}

	a.size = a.terminal.Size()
	if err := a.terminal.Clear(); err != nil {
		return err
	}

	return a.controller.Redraw()
}

func (a *App) Run() error {
	defer a.terminal.Close()
	defer a.controller.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		// For every ticker, we want to check if the terminal should be redraw
		select {
		case <-ticker.C:
			if err := a.resize(); err != nil {
				return err
			}
			if a.ctx.Err() != nil {
				return nil
			}
		case <-a.redrawCh:
			if err := a.controller.Redraw(); err != nil {
				return fmt.Errorf("termdash.controller.Redraw => %v", err)
			}

		case err := <-a.errorCh:
			if err != nil {
				return err
			}
		case <-a.ctx.Done():
			return nil
		}
	}
}
