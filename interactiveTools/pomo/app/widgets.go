package app

import (
	"context"

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

func (w *widgets) updateWidgets(redrawCh chan<- bool, txtInfo, txtType string) {
	if txtInfo != "" {
		w.updateTxtInfo <- txtInfo
	}

	if txtType != "" {
		w.updateTxtType <- txtType
	}
	redrawCh <- true
}

func newWidget(ctx context.Context, errorCh chan<- error) (*widgets, error) {
	w := &widgets{}
	var err error
	w.updateDonTimer = make(chan []int)
	w.updateTxtType = make(chan string)
	w.updateTxtInfo = make(chan string)
	w.updateTxtTimer = make(chan string)

	w.displayType, err = newSegmentDisplay(ctx, w.updateTxtType, errorCh)
	if err != nil {
		return nil, err
	}
	w.donutTimer, err = newDonut(ctx)
	if err != nil {
		return nil, err
	}

	w.txtInfo, err = newText(ctx, w.updateTxtInfo, errorCh)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func newSegmentDisplay(ctx context.Context, updateText <-chan string, errorCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}
	// Goroutine to update SegmentDisplay
	go func() {
		for {
			select {
			case t := <-updateText:
				if t == "" {
					t = " "
				}

				errorCh <- sd.Write([]*segmentdisplay.TextChunk{
					segmentdisplay.NewChunk(t),
				})
			case <-ctx.Done():
				return
			}
		}
	}()

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

	if err := d.Percent(progress); err != nil {
		return nil, err
	}
	return d, nil
}

func newText(
	ctx context.Context, updateText <-chan string, errorCh chan<- error,
) (*text.Text, error) {
	txt, err := text.New()
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case t := <-updateText:
				txt.Reset()
				errorCh <- txt.Write(t)
			case <-ctx.Done():
				return
			}
		}
	}()
	return txt, nil
}
