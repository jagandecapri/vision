package data

import (
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"html/template"
	"bytes"
	"github.com/jagandecapri/vision/preprocess"
	"github.com/jagandecapri/vision/tree"
	"encoding/json"
	"fmt"
	"strings"
)

type SQL struct{
	db_name string
	metadata_table string
	agg_src_table string
	agg_dst_table string
	agg_srcdst_table string
	flow_keys map[string]string
	delta_t string
	time string
}

func (s *SQL) SetupDb(){
	db, err := sql.Open("sqlite3", s.db_name)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//CREATE TABLE TO STORE METADATA
	metadata_table := `CREATE TABLE IF NOT EXISTS ` +s.metadata_table+ ` (
		id	INTEGER PRIMARY KEY AUTOINCREMENT,
		delta_t	TEXT,
		table_type TEXT,
		table_name TEXT,
		time NUMERIC,
		UNIQUE(delta_t, table_type))`

	_, err = db.Exec(metadata_table)
	if err != nil {
		log.Printf("%q: %s\n", err, metadata_table)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("REPLACE INTO metadata(delta_t, table_type, table_name, time) values(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(s.delta_t, "agg_src", s.agg_src_table, s.time); err != nil {
		log.Fatal(err)
	}
	if _, err := stmt.Exec(s.delta_t, "agg_dst", s.agg_dst_table, s.time); err != nil {
		log.Fatal(err)
	}
	if _, err := stmt.Exec(s.delta_t, "agg_srcdst", s.agg_srcdst_table, s.time); err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	//DELETE TABLES


	//CREATE TABLES TO STORE DATA
	t := template.New("create tables")
	t, _ = t.Parse(`CREATE TABLE IF NOT EXISTS ` + "`{{.TableName}}`" + `(
		id	INTEGER PRIMARY KEY AUTOINCREMENT,
		batch INTEGER,
		flow_key TEXT,
		agg_src TEXT,
		agg_dst TEXT,
		nbSrcPort REAL,
		nbDstPort REAL,
		nbSrcs REAL,
		nbDsts REAL,
		perSYN REAL,
		perACK REAL,
		perICMP REAL,
		perRST REAL,
		perFIN REAL,
		perCWR REAL,
		perURG REAL,
		avgPktSize REAL,
		meanTTL REAL,
		last_packet_timestamp	NUMERIC
	)`)

	var tpl bytes.Buffer
	var result string

	if err := t.Execute(&tpl, struct{TableName string}{TableName: s.agg_src_table}); err != nil {
		log.Println(err)
	}

	result = tpl.String()

	_, err = db.Exec(result)
	if err != nil {
		log.Printf("agg_src: %q: %s\n", err, result)
		return
	}

	tpl.Reset()
	if err := t.Execute(&tpl, struct{TableName string}{TableName: s.agg_dst_table}); err != nil {
		log.Println(err)
	}

	result = tpl.String()

	_, err = db.Exec(result)
	if err != nil {
		log.Printf("agg_dst: %q: %s\n", err, result)
		return
	}

	tpl.Reset()
	if err := t.Execute(&tpl, struct{TableName string}{TableName: s.agg_srcdst_table}); err != nil {
		log.Println(err)
	}

	result = tpl.String()

	_, err = db.Exec(result)
	if err != nil {
		log.Printf("agg_srcdst %q: %s\n", err, result)
		return
	}
}

func (s *SQL) WriteToDb(acc_c preprocess.AccumulatorChannels, done chan struct{}){
	//CREATE TABLES TO STORE DATA
	t := template.New("insert data to tables")
	t, _ = t.Parse(`INSERT INTO ` + "`{{.TableName}}`" +`(batch, 
						flow_key, 
						agg_src, 
						agg_dst, 
						last_packet_timestamp,
						nbSrcPort, 
						nbDstPort, 
						nbSrcs, 
						nbDsts, 
						perSYN, 
						perACK,
						perICMP, 
						perRST, 
						perFIN, 
						perCWR, 
						perURG,
						avgPktSize, 
						meanTTL) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)


	var tpl bytes.Buffer
	var result string
	batch_counter_agg_src := 1
	batch_counter_agg_dst := 1
	batch_counter_agg_srcdst := 1

	go func(){
		var tmp []tree.Point

		log.Printf("*")

		for{
			select{
				case X := <-acc_c.AggSrc:
					tmp = X

					tpl.Reset()
					if err := t.Execute(&tpl, struct{TableName string}{TableName: s.agg_src_table}); err != nil {
						log.Fatal(err)
					}

					db, err := sql.Open("sqlite3", s.db_name)
					if err != nil {
						log.Fatal(err)
					}

					tx, err := db.Begin()
					if err != nil {
						log.Fatal(err)
					}

					result = tpl.String()

					stmt, err := tx.Prepare(result)
					if err != nil {
						log.Fatal(err)
					}

					for _, p := range tmp {
						var json_string []byte
						var err error
						if json_string, err = json.Marshal(p.Key); err != nil{
							log.Fatal(err)
						}

						var tmp1 = make([]interface{}, 18)

						tmp1[0] = batch_counter_agg_src
						tmp1[1] = string(json_string)
						tmp1[2] = strings.Join(p.Key.SrcIP, ",")
						tmp1[3] = strings.Join(p.Key.DstIP, ",")
						tmp1[4] = nil
						tmp1[5] = p.Vec_map["nbSrcPort"]
						tmp1[6]	= p.Vec_map["nbDstPort"]
						tmp1[7]	= p.Vec_map["nbSrcs"]
						tmp1[8]	= p.Vec_map["nbDsts"]
						tmp1[9]	= p.Vec_map["perSYN"]
						tmp1[10] = p.Vec_map["perACK"]
						tmp1[11] = p.Vec_map["perICMP"]
						tmp1[12] = p.Vec_map["perRST"]
						tmp1[13] = p.Vec_map["perFIN"]
						tmp1[14] = p.Vec_map["perCWR"]
						tmp1[15] = p.Vec_map["perURG"]
						tmp1[16] = p.Vec_map["avgPktSize"]
						tmp1[17] = p.Vec_map["meanTTL"]

						_, err = stmt.Exec(tmp1...)
						if err != nil {
							log.Fatal(tmp1, err)
						}
					}
					tx.Commit()
					stmt.Close()
					db.Close()
					batch_counter_agg_src++
			case X := <-acc_c.AggDst:
					tmp = X

					tpl.Reset()
					if err := t.Execute(&tpl, struct{TableName string}{TableName: s.agg_dst_table}); err != nil {
						log.Fatal(err)
					}

					db, err := sql.Open("sqlite3", s.db_name)
					if err != nil {
						log.Fatal(err)
					}

					tx, err := db.Begin()
					if err != nil {
						log.Fatal(err)
					}

					result = tpl.String()

					stmt, err := tx.Prepare(result)
					if err != nil {
						log.Fatal(err)
					}

					for _, p := range tmp {
						var json_string []byte
						var err error
						if json_string, err = json.Marshal(p.Key); err != nil{
							log.Fatal(err)
						}

						var tmp1 = make([]interface{}, 18)

						tmp1[0] = batch_counter_agg_src
						tmp1[1] = string(json_string)
						tmp1[2] = strings.Join(p.Key.SrcIP, ",")
						tmp1[3] = strings.Join(p.Key.DstIP, ",")
						tmp1[4] = nil
						tmp1[5] = p.Vec_map["nbSrcPort"]
						tmp1[6]	= p.Vec_map["nbDstPort"]
						tmp1[7]	= p.Vec_map["nbSrcs"]
						tmp1[8]	= p.Vec_map["nbDsts"]
						tmp1[9]	= p.Vec_map["perSYN"]
						tmp1[10] = p.Vec_map["perACK"]
						tmp1[11] = p.Vec_map["perICMP"]
						tmp1[12] = p.Vec_map["perRST"]
						tmp1[13] = p.Vec_map["perFIN"]
						tmp1[14] = p.Vec_map["perCWR"]
						tmp1[15] = p.Vec_map["perURG"]
						tmp1[16] = p.Vec_map["avgPktSize"]
						tmp1[17] = p.Vec_map["meanTTL"]

						_, err = stmt.Exec(tmp1...)
						if err != nil {
							log.Fatal(tmp1, err)
						}
					}
					tx.Commit()
					stmt.Close()
					db.Close()
					batch_counter_agg_dst++
				case X := <-acc_c.AggSrcDst:
					tmp = X

					tpl.Reset()
					if err := t.Execute(&tpl, struct{TableName string}{TableName: s.agg_srcdst_table}); err != nil {
						log.Fatal(err)
					}

					db, err := sql.Open("sqlite3", s.db_name)
					if err != nil {
						log.Fatal(err)
					}

					tx, err := db.Begin()
					if err != nil {
						log.Fatal(err)
					}

					result = tpl.String()

					stmt, err := tx.Prepare(result)
					if err != nil {
						log.Fatal(err)
					}

					for _, p := range tmp {
						var json_string []byte
						var err error
						if json_string, err = json.Marshal(p.Key); err != nil{
							log.Fatal(err)
						}

						var tmp1 = make([]interface{}, 18)

						tmp1[0] = batch_counter_agg_src
						tmp1[1] = string(json_string)
						tmp1[2] = strings.Join(p.Key.SrcIP, ",")
						tmp1[3] = strings.Join(p.Key.DstIP, ",")
						tmp1[4] = nil
						tmp1[5] = p.Vec_map["nbSrcPort"]
						tmp1[6]	= p.Vec_map["nbDstPort"]
						tmp1[7]	= p.Vec_map["nbSrcs"]
						tmp1[8]	= p.Vec_map["nbDsts"]
						tmp1[9]	= p.Vec_map["perSYN"]
						tmp1[10] = p.Vec_map["perACK"]
						tmp1[11] = p.Vec_map["perICMP"]
						tmp1[12] = p.Vec_map["perRST"]
						tmp1[13] = p.Vec_map["perFIN"]
						tmp1[14] = p.Vec_map["perCWR"]
						tmp1[15] = p.Vec_map["perURG"]
						tmp1[16] = p.Vec_map["avgPktSize"]
						tmp1[17] = p.Vec_map["meanTTL"]

						_, err = stmt.Exec(tmp1...)
						if err != nil {
							log.Fatal(tmp1, err)
						}
					}
					tx.Commit()
					stmt.Close()
					db.Close()
					batch_counter_agg_srcdst++
				case <-done:
					return
				default:
			}
		}
	}()
}

func NewSQL(acc_c preprocess.AccumulatorChannels, done chan struct{}, delta_t time.Duration) SQL{
	now := time.Now()
	now_string := now.Format(time.RFC3339)
	fmt.Println(now_string)
	sql := SQL{
		db_name: "./vision.db",
		metadata_table: "metadata",
		agg_src_table: "agg_src_" + now_string,
		agg_dst_table: "agg_dst_" + now_string,
		agg_srcdst_table: "agg_srcdst_" + now_string,
		delta_t: delta_t.String(),
		time: now_string,
	}
	sql.SetupDb()
	sql.WriteToDb(acc_c, done)

	return sql
}
