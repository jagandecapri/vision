package data

import (
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jagandecapri/vision/preprocess"
	"os"
	"time"
	"fmt"
	"database/sql"
)

var delta_t = 300 * time.Millisecond
var window = 15 * time.Second
var WINDOW_ARR_LEN = int(window.Seconds()/delta_t.Seconds())
var Point_ctr = 0

func WriteToDb(acc_c preprocess.AccumulatorChannel){
	os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	create_agg_src := `CREATE TABLE IF NOT EXISTS agg_src (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
	flow_data	TEXT,
	last_packet_timestamp	NUMERIC
);`
	_, err = db.Exec(create_agg_src)
	if err != nil {
		log.Printf("%q: %s\n", err, create_agg_src)
		return
	}

	create_agg_dst := `CREATE TABLE IF NOT EXISTS agg_src (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
	flow_data	TEXT,
	last_packet_timestamp	NUMERIC
);`
	_, err = db.Exec(create_agg_dst)
	if err != nil {
		log.Printf("%q: %s\n", err, create_agg_dst)
		return
	}

	create_agg_srcdst := `CREATE TABLE IF NOT EXISTS agg_src (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
	flow_data	TEXT,
	last_packet_timestamp	NUMERIC
);`
	_, err = db.Exec(create_agg_src)
	if err != nil {
		log.Printf("%q: %s\n", err, create_agg_srcdst)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func WindowTimeSlide(ch chan preprocess.PacketData, acc_c preprocess.AccumulatorChannels, done chan struct{}){

	go func(){
		acc := preprocess.NewAccumulator()
		time_init := time.Now()
		time_counter := time.Time{}

		for{
			select{
			case pd := <- ch:
				packet_time := pd.Metadata.Timestamp

				if time_counter.IsZero(){
					fmt.Println("Initialize Time")
					time_counter = packet_time
					acc = preprocess.NewAccumulator()
				}

				if !time_counter.IsZero() && packet_time.After(time_counter.Add(delta_t)){
					fmt.Println("packet_time > time_counter")
					X := acc.GetMicroSlot()
					acc_c.AggSrc <- X.AggSrc
					acc_c.AggDst <- X.AggDst
					acc_c.AggSrcDst <- X.AggSrcDst
					log.Println("Time to read data:", time.Since(time_init))
					time_init = time.Now()
					time_counter = time.Time{}
				}

				acc.AddPacket(pd.Data)
			case <-done:
				return
			default:
			}
		}
	}()
}

func main(){
	ch := make(chan preprocess.PacketData)
	done := make(chan struct{})

	pcap_file_path := os.Getenv("PCAP_FILE")
	if pcap_file_path == ""{
		pcap_file_path = "C:\\Users\\Jack\\Downloads\\201705021400.pcap"
	}

	handleRead, err := pcap.OpenOffline(pcap_file_path)

	if(err != nil){
		log.Fatal(err)
	}

	for {
		data, ci, err := handleRead.ReadPacketData()
		if err != nil && err != io.EOF {
			close(done)
			log.Fatal(err)
		} else if err == io.EOF {
			close(done)
			break
		} else {
			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
			ch <- preprocess.PacketData{Data: packet, Metadata: ci}
		}
	}
}
