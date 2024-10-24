package app

import (
	"context"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
)

type widgets struct {
	donTimer       *donut.Donut
	disType        *segmentdisplay.SegmentDisplay
	txtInfo        *text.Text
	txtTimer       *text.Text
	updateTxtInfo  chan string
	updateTxtType  chan string
	updateTxtTime  chan string
	updateDonTimer chan []int
}

func (w *widgets) update(timer []int, txtType, txtInfo, txtTimer string, redrawCh chan<- bool) {
	if txtInfo != "" {
		w.updateTxtInfo <- txtInfo
	}

	if txtType != "" {
		w.updateTxtType <- txtType
	}

	if txtTimer != "" {
		w.updateTxtTime <- txtTimer
	}
	if len(timer) > 0 {
		w.updateDonTimer <- timer
	}
	redrawCh <- true
}

func newWidgets(ctx context.Context, errorCh chan<- error) (*widgets, error) {
	w := &widgets{}
	var err error

	w.updateDonTimer = make(chan []int)
	w.updateTxtType = make(chan string)
	w.updateTxtInfo = make(chan string)
	w.updateTxtTime = make(chan string)

	// w.donTimer, err = newDonut(ctx, w.updateDonTimer, errorCh)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// w.txtTimer , err = newText(ctx, w.updateTxtType)
}

func newText(ctx context.Context, updateText <-chan string, errorCh chan<- error) (*text.Text, error) {
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

func newDonut(ctx context.Context, donUpdater <-chan []int, errorCh chan<- error) (*donut.Donut, error) {
	donut, err := donut.New(
		donut.Clockwise(),
		donut.CellOpts(cell.FgColor(cell.ColorBlue)),
	)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case d := <-donUpdater:
				if d[0] < d[1] {
					errorCh <- donut.Absolute(d[0], d[1])
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return donut, nil
}

func newSegmentDisplay(ctx context.Context, updateText <-chan string, errorCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}
	// Go routine to update the segment display
	go func() {
		for {
			select {
			case t := <-updateText:
				if t == "" {
					t = " "
				}
				errorCh <- sd.Write([]*segmentdisplay.TextChunk{segmentdisplay.NewChunk(t)})
			case <-ctx.Done():
				return
			}
		}
	}()

	return sd, nil
}
