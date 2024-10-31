package app

import (
	"fmt"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/button"
)

type buttonSet struct {
	btStart *button.Button
	btPause *button.Button
}

func newButtons() (*buttonSet, error) {
	startB, err := button.New("(s)tart", func() error {
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
