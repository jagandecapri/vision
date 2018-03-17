package data

import (
	"testing"
	"github.com/jagandecapri/vision/preprocess"
	"time"
	"log"
)

func TestRun(t *testing.T) {
	Run()
}

func TestNewSQLRead(t *testing.T) {
	delta_t := 300 * time.Millisecond
	sql := NewSQLRead("201705021400", delta_t)

	acc_c := preprocess.AccumulatorChannels{
		AggSrc: make(preprocess.AccumulatorChannel),
		AggDst: make(preprocess.AccumulatorChannel),
		AggSrcDst: make(preprocess.AccumulatorChannel),
	}

	done := sql.ReadFromDb(acc_c)

	now_src := time.Now()
	now_dst := time.Now()
	now_srcdst := time.Now()

	for{
		select{
			case <-acc_c.AggSrc:
				tmp := time.Since(now_src)
				log.Println("Aggsrc data received in ", tmp)
				now_src = time.Now()
			case <-acc_c.AggDst:
				tmp := time.Since(now_dst)
				log.Println("Aggdst data received in ", tmp)
				now_dst = time.Now()
			case <-acc_c.AggSrcDst:
				tmp := time.Since(now_srcdst)
				log.Println("Aggsrcdst data received in ", tmp)
				now_srcdst = time.Now()
			case <-done:
				return
		}
	}
}