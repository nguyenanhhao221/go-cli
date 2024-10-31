package app

import (
	"context"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

type App struct {
	ctx        context.Context
	controller *termdash.Controller
	terminal   *tcell.Terminal
}

func New() (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

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

	c, err := newGrid(ctx, term)
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
	}, nil
}

func (a *App) Run() error {
	// TODO temporary to prevent error when quit need to find a better way
	defer func() {
		time.Sleep(100 * time.Millisecond) // Allow Termdash to finish processing
		a.controller.Close()
		a.terminal.Close()
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if a.ctx.Err() != nil { // Exit loop if context is done
				return nil
			}
			if err := a.controller.Redraw(); err != nil {
				return err
			}
		case <-a.ctx.Done():
			return nil
		}
	}
}
