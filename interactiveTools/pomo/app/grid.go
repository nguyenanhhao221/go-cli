package app

import (
	"context"

	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

// newGrid Get the container and define the layout for widgets
func newGrid(ctx context.Context, t terminalapi.Terminal) (*container.Container, error) {
	widgets, err := newWidget(ctx)
	if err != nil {
		return nil, err
	}
	b, err := newButtons()
	if err != nil {
		return nil, err
	}

	builder := grid.New()

	// First row
	builder.Add(
		grid.RowHeightPerc(70,
			grid.ColWidthPercWithOpts(
				60,
				[]container.Option{
					container.AlignHorizontal(align.HorizontalCenter),
				},
				grid.RowHeightPerc(25, grid.Widget(widgets.donutTimer,
					container.BorderTitle("Press Q to quit"),
					container.Border(linestyle.Light))),
			),
			grid.ColWidthPercWithOpts(
				40,
				[]container.Option{
					container.AlignHorizontal(align.HorizontalCenter),
				},
				// The pomodoro text section
				grid.RowHeightPerc(70, grid.Widget(widgets.displayType,
					container.Border(linestyle.Light))),
				// The message
				grid.RowHeightPerc(25, grid.Widget(widgets.displayType,
					container.AlignHorizontal(align.HorizontalCenter),
					container.Border(linestyle.Light))),
			),
		),
	)

	// Second row
	builder.Add(
		grid.RowHeightPerc(30,
			grid.ColWidthPerc(50, grid.Widget(b.btStart)),
			grid.ColWidthPerc(50, grid.Widget(b.btPause))),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	c, err := container.New(t, gridOpts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}
