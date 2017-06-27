package utils

import (
	"github.com/jagandecapri/vision/orunada/grid"
)

func Broadcaster(bcast *ThreadSafeSlice, ch chan grid.HttpData) {
	for {
		msg := <-ch
		bcast.Iter(func(w *Worker) { w.Source <- msg })
	}
}