package app

import (
	"context"
	"fmt"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/button"
	"haonguyen.tech/interactiveTools/pomo/pomodoro"
)

type buttonSet struct {
	btStart *button.Button
	btPause *button.Button
}

// newButtons Display the buttons, also include callback when a button is pressing to control the UI
func newButtons(ctx context.Context, config *pomodoro.IntervalConfig,
	w *widgets, errorCh chan error, redrawCh chan<- bool,
) (*buttonSet, error) {
	startInterval := func() {
		i, err := pomodoro.GetInterval(config)
		errorCh <- err

		start := func(i pomodoro.Interval) {
			message := "Take a break"
			if i.Category == pomodoro.CategoryPomodoro {
				message = "Focus on your task"
			}

			w.updateWidgets(redrawCh, message, i.Category, "", []int{})
		}

		periodic := func(i pomodoro.Interval) {
			w.updateWidgets(redrawCh, "", "", fmt.Sprint(i.PlannedDuration-i.ActualDuration), []int{int(i.ActualDuration), int(i.PlannedDuration)})
		}

		end := func(pomodoro.Interval) {
			w.updateWidgets(redrawCh, "Nothing running...", i.Category, "", []int{})
		}

		errorCh <- i.Start(ctx, config, start, periodic, end)
	}
	startB, err := button.New("(s)tart", func() error {
		go startInterval()
		return nil
	},
		button.GlobalKey('s'),
		button.WidthFor("(s)tart"),
	)
	if err != nil {
		return nil, fmt.Errorf("newButtons: %w", err)
	}

	pauseB, err := button.New("(p)ause", func() error {
		return nil
	},
		button.GlobalKey('p'),
		button.WidthFor("(p)ause"),
		button.FillColor(cell.ColorNumber(220)),
	)
	if err != nil {
		return nil, fmt.Errorf("newButtons: %w", err)
	}

	return &buttonSet{btStart: startB, btPause: pauseB}, nil
}
