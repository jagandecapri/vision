package data

import (
	"testing"
	"github.com/jagandecapri/vision/preprocess"
	"time"
	"log"
)

func TestRunt(t *testing.T) {
	Run()
}

func TestNewSQLRead(t *testing.T) {
	delta_t := 300 * time.Millisecond
	sql := NewSQLRead(delta_t)

	acc_c := preprocess.AccumulatorChannels{
		AggSrc: make(preprocess.AccumulatorChannel),
		AggDst: make(preprocess.AccumulatorChannel),
		AggSrcDst: make(preprocess.AccumulatorChannel),
	}

	done := sql.ReadFromDb(acc_c)

	for{
		select{
			case <-acc_c.AggSrc:
				log.Println("Received aggsrc data")
			case <-acc_c.AggDst:
				log.Println("Received aggdst data")
			case <-acc_c.AggSrcDst:
				log.Println("Received aggsrcdst data")
			case <-done:
				return
		}
	}
}