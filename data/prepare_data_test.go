package data

import (
	"testing"
	"github.com/jagandecapri/vision/preprocess"
	"time"
	"log"
	"flag"
	"os"
)

var pcap_file_path = flag.String("pcap_file_path", "", "pcap_file_path")
var db_name = flag.String("db_name", "", "db_name")
var delta_t = flag.Duration("delta_t", 300 * time.Millisecond, "Delta time")

func TestRun(t *testing.T) {
	flag.Parse()

	if *pcap_file_path == "" || *db_name == ""{
		flag.PrintDefaults()
		os.Exit(1)
	}

	Run(*pcap_file_path, *db_name, *delta_t)
}

func TestNewSQLRead(t *testing.T) {
	delta_t := 300 * time.Millisecond
	sql := NewSQLRead("201705021400", delta_t)
	done := make(chan struct{})

	acc_c := preprocess.AccumulatorChannels{
		AggSrc: make(preprocess.AccumulatorChannel),
		AggDst: make(preprocess.AccumulatorChannel),
		AggSrcDst: make(preprocess.AccumulatorChannel),
	}

	sql.ReadFromDb(acc_c)

	now_src := time.Now()
	now_dst := time.Now()
	now_srcdst := time.Now()

	counter_done := 0

	for{
		select{
			case _, ok := <-acc_c.AggSrc:
				if ok{
					tmp := time.Since(now_src)
					log.Println("Aggsrc data received in ", tmp)
					now_src = time.Now()
				} else {
					counter_done++
				}
			case _, ok := <-acc_c.AggDst:
				if ok{
					tmp := time.Since(now_dst)
					log.Println("Aggdst data received in ", tmp)
					now_dst = time.Now()
				} else {
					counter_done++
				}
			case _, ok := <-acc_c.AggSrcDst:
				if ok{
					tmp := time.Since(now_srcdst)
					log.Println("Aggsrcdst data received in ", tmp)
					now_srcdst = time.Now()
				} else {
					counter_done++
				}
			case <-done:
				return
		}

		if counter_done == 3{
			break
		}
	}
}