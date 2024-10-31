package app

import (
	"context"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
)

type widgets struct {
	donutTimer     *donut.Donut
	displayType    *segmentdisplay.SegmentDisplay
	txtInfo        *text.Text
	txtTimer       *text.Text
	updateDonTimer chan []int
	updateTxtInfo  chan string
	updateTxtTimer chan string
	updateTxtType  chan string
}

func newWidget(ctx context.Context) (*widgets, error) {
	w := &widgets{}
	var err error
	// updateText := make(chan string)
	// errorCh := make(chan error)
	w.displayType, err = newSegmentDisplay()
	if err != nil {
		return nil, err
	}
	w.donutTimer, err = newDonut(ctx)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func newSegmentDisplay() (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}
	//TODO: update once the implementation for pomodoro is finished
	// Goroutine to update SegmentDisplay
	// go func() {
	// 	for {
	// 		select {
	// 		case t := <-updateText:
	// 			if t == "" {
	// 				t = " "
	// 			}
	//
	// 			errorCh <- sd.Write([]*segmentdisplay.TextChunk{
	// 				segmentdisplay.NewChunk(t),
	// 			})
	// 		case <-ctx.Done():
	// 			return
	// 		}
	// 	}
	// }()

	t := "Pomodoro"
	if err := sd.Write([]*segmentdisplay.TextChunk{
		segmentdisplay.NewChunk(t),
	}); err != nil {
		return nil, err
	}
	return sd, nil
}

func newDonut(ctx context.Context) (*donut.Donut, error) {
	d, err := donut.New(donut.CellOpts(
		cell.FgColor(cell.ColorNumber(33))),
	)
	if err != nil {
		return nil, err
	}

	const start = 35
	progress := start

	go periodic(ctx, 500*time.Millisecond, func() error {
		if err := d.Percent(progress); err != nil {
			return err
		}
		progress++
		if progress > 100 {
			progress = start
		}
		return nil
	})
	return d, nil
}

func periodic(ctx context.Context, interval time.Duration, fn func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := fn(); err != nil {
				panic(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
