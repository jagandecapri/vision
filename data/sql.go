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
	"strings"
	"sync"
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
		created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now', 'localtime')),
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
		nbPacket REAL,
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
		last_packet_timestamp	NUMERIC,
		created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now', 'localtime'))
	)`)

	var tpl bytes.Buffer
	var result string

	if err := t.Execute(&tpl, struct{TableName template.HTML}{TableName: template.HTML(s.agg_src_table)}); err != nil {
		log.Println(err)
	}

	result = tpl.String()

	_, err = db.Exec(result)
	if err != nil {
		log.Printf("agg_src: %q: %s\n", err, result)
		return
	}

	tpl.Reset()
	if err := t.Execute(&tpl, struct{TableName template.HTML}{TableName: template.HTML(s.agg_dst_table)}); err != nil {
		log.Println(err)
	}

	result = tpl.String()

	_, err = db.Exec(result)
	if err != nil {
		log.Printf("agg_dst: %q: %s\n", err, result)
		return
	}

	tpl.Reset()
	if err := t.Execute(&tpl, struct{TableName template.HTML}{TableName: template.HTML(s.agg_srcdst_table)}); err != nil {
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
						nbPacket,
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
						meanTTL) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)


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
					log.Println("Writing aggsrc data")
					tmp = X

					tpl.Reset()
					if err := t.Execute(&tpl, struct{TableName template.HTML}{TableName: template.HTML(s.agg_src_table)}); err != nil {
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

						var tmp1 = make([]interface{}, 19)

						tmp1[0] = batch_counter_agg_src
						tmp1[1] = string(json_string)
						tmp1[2] = strings.Join(p.Key.SrcIP, ",")
						tmp1[3] = strings.Join(p.Key.DstIP, ",")
						tmp1[4] = nil
						tmp1[5] = p.Vec_map["nbPacket"]
						tmp1[6] = p.Vec_map["nbSrcPort"]
						tmp1[7]	= p.Vec_map["nbDstPort"]
						tmp1[8]	= p.Vec_map["nbSrcs"]
						tmp1[9]	= p.Vec_map["nbDsts"]
						tmp1[10] = p.Vec_map["perSYN"]
						tmp1[11] = p.Vec_map["perACK"]
						tmp1[12] = p.Vec_map["perICMP"]
						tmp1[13] = p.Vec_map["perRST"]
						tmp1[14] = p.Vec_map["perFIN"]
						tmp1[15] = p.Vec_map["perCWR"]
						tmp1[16] = p.Vec_map["perURG"]
						tmp1[17] = p.Vec_map["avgPktSize"]
						tmp1[18] = p.Vec_map["meanTTL"]

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
					log.Println("Writing aggdst data")
					tmp = X

					tpl.Reset()
					if err := t.Execute(&tpl, struct{TableName template.HTML}{TableName: template.HTML(s.agg_dst_table)}); err != nil {
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

						var tmp1 = make([]interface{}, 19)

						tmp1[0] = batch_counter_agg_dst
						tmp1[1] = string(json_string)
						tmp1[2] = strings.Join(p.Key.SrcIP, ",")
						tmp1[3] = strings.Join(p.Key.DstIP, ",")
						tmp1[4] = nil
						tmp1[5] = p.Vec_map["nbPacket"]
						tmp1[6] = p.Vec_map["nbSrcPort"]
						tmp1[7]	= p.Vec_map["nbDstPort"]
						tmp1[8]	= p.Vec_map["nbSrcs"]
						tmp1[9]	= p.Vec_map["nbDsts"]
						tmp1[10] = p.Vec_map["perSYN"]
						tmp1[11] = p.Vec_map["perACK"]
						tmp1[12] = p.Vec_map["perICMP"]
						tmp1[13] = p.Vec_map["perRST"]
						tmp1[14] = p.Vec_map["perFIN"]
						tmp1[15] = p.Vec_map["perCWR"]
						tmp1[16] = p.Vec_map["perURG"]
						tmp1[17] = p.Vec_map["avgPktSize"]
						tmp1[18] = p.Vec_map["meanTTL"]

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
					log.Println("Writing aggsrcdst data")
					tmp = X

					tpl.Reset()
					if err := t.Execute(&tpl, struct{TableName template.HTML}{TableName: template.HTML(s.agg_srcdst_table)}); err != nil {
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

						var tmp1 = make([]interface{}, 19)

						tmp1[0] = batch_counter_agg_srcdst
						tmp1[1] = string(json_string)
						tmp1[2] = strings.Join(p.Key.SrcIP, ",")
						tmp1[3] = strings.Join(p.Key.DstIP, ",")
						tmp1[4] = nil
						tmp1[5] = p.Vec_map["nbPacket"]
						tmp1[6] = p.Vec_map["nbSrcPort"]
						tmp1[7]	= p.Vec_map["nbDstPort"]
						tmp1[8]	= p.Vec_map["nbSrcs"]
						tmp1[9]	= p.Vec_map["nbDsts"]
						tmp1[10] = p.Vec_map["perSYN"]
						tmp1[11] = p.Vec_map["perACK"]
						tmp1[12] = p.Vec_map["perICMP"]
						tmp1[13] = p.Vec_map["perRST"]
						tmp1[14] = p.Vec_map["perFIN"]
						tmp1[15] = p.Vec_map["perCWR"]
						tmp1[16] = p.Vec_map["perURG"]
						tmp1[17] = p.Vec_map["avgPktSize"]
						tmp1[18] = p.Vec_map["meanTTL"]

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

func (s *SQL) ReadFromDb(acc_c preprocess.AccumulatorChannels){
	var batch_counter_agg_src, batch_counter_agg_dst, batch_counter_agg_srcdst int

	log.Println(s.db_name)
	db, err := sql.Open("sqlite3", s.db_name)
	if err != nil {
		log.Fatal(err)
	}

	var query string

	query = "SELECT MAX(batch) FROM `" + s.agg_src_table + "`"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(query, err)
	}

	for rows.Next() {
		err = rows.Scan(&batch_counter_agg_src)
		if err != nil {
			log.Fatal(query, err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(query, err)
	}

	rows.Close()

	query = "SELECT MAX(batch) FROM `" + s.agg_dst_table + "`"
	rows, err = db.Query(query)
	if err != nil {
		log.Fatal(query, err)
	}

	for rows.Next() {
		err = rows.Scan(&batch_counter_agg_dst)
		if err != nil {
			log.Fatal(query, err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(query, err)
	}

	rows.Close()

	query = "SELECT MAX(batch) FROM `" + s.agg_srcdst_table + "`"
	rows, err = db.Query(query)
	if err != nil {
		log.Fatal(query, err)
	}

	for rows.Next() {
		err = rows.Scan(&batch_counter_agg_srcdst)
		if err != nil {
			log.Fatal(query, err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(query, err)
	}

	rows.Close()

	db.Close()

	go func(){
		var wg sync.WaitGroup
		wg.Add(3)

		read_fn := func(acc_c preprocess.AccumulatorChannel, table string, wg *sync.WaitGroup){
				defer func(){
					log.Println("close windowtimeslider channel")
					close(acc_c)
					wg.Done()
				}()

				t := template.New("select data from agg_src")
				t, err = t.Parse(`SELECT id, flow_key, nbPacket, nbSrcPort,
			nbDstPort, nbSrcs, nbDsts, perSYN, perACK, perICMP, perRST, perFIN, perCWR, perURG,
			avgPktSize, meanTTL FROM ` + "`{{.TableName}}`" + " WHERE batch={{.Batch}}" )
				if err != nil{
					log.Fatal(err)
				}

				var buf bytes.Buffer
				var query string

				for i := 1; i <= batch_counter_agg_src; i++ {
					log.Println("Read loop: ", i)
					tpl_data := struct{TableName template.HTML
						Batch int}{TableName: template.HTML(table), Batch: i}
					if err := t.Execute(&buf, tpl_data); err != nil {
						log.Fatal("batch looping src err: ", err)
					}

					query = buf.String()
					points := s.IterateRows(query)
					acc_c <- points
					buf.Reset()
				}
				return
		}

		go read_fn(acc_c.AggSrc, s.agg_src_table, &wg)
		go read_fn(acc_c.AggDst, s.agg_dst_table, &wg)
		go read_fn(acc_c.AggSrcDst, s.agg_srcdst_table, &wg)

		wg.Wait()
		return
	}()
}


func (s SQL) IterateRows(query string) []tree.Point{
	db, err := sql.Open("sqlite3", s.db_name)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	points := []tree.Point{}

	for rows.Next() {
		var id int
		var flow_key sql.RawBytes
		var nbPacket, nbSrcPort, nbDstPort,
			nbSrcs, nbDsts, perSYN,
			perACK, perICMP, perRST,
			perFIN, perCWR, perURG,
			avgPktSize, meanTTL float64

		err = rows.Scan(&id, &flow_key, &nbPacket, &nbSrcPort,
			&nbDstPort, &nbSrcs, &nbDsts,
				&perSYN, &perACK, &perICMP,
					&perRST, &perFIN, &perCWR,
						&perURG, &avgPktSize, &meanTTL)

		if err != nil {
			log.Fatal(err)
		}

		x := map[string]float64{
			"nbPacket": nbPacket,
			"nbSrcPort": nbSrcPort,
			"nbDstPort": nbDstPort,
			"nbSrcs": nbSrcs,
			"nbDsts": nbDsts,
			"perSYN": perSYN,
			"perACK": perACK,
			"perICMP": perICMP,
			"perRST": perRST,
			"perFIN": perFIN,
			"perCWR": perCWR,
			"perURG": perURG,
			"avgPktSize": avgPktSize,
			"meanTTL": meanTTL,
		}

		point_key := tree.PointKey{}
		json.Unmarshal(flow_key, &point_key)

		pnt := tree.Point{Id: id, Key: point_key, Vec_map: x}
		points = append(points, pnt)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()

	return points
}

func NewSQL(db_name string, acc_c preprocess.AccumulatorChannels, done chan struct{}, delta_t time.Duration) SQL{
	now := time.Now()
	now_string := now.Format(time.RFC3339)
	sql := SQL{
		db_name: db_name,
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

func NewSQLRead(db_name string, delta_t time.Duration) SQL{
	metadata_table := "metadata"
	var agg_src_table, agg_dst_table, agg_srcdst_table string

	db, err := sql.Open("sqlite3", db_name)
	if err != nil {
		log.Fatal("Open database err: ", err)
	}
	defer db.Close()

	query := "SELECT table_name, table_type FROM " +metadata_table + " WHERE delta_t = '" + delta_t.String() + "'"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var table_name string
		var table_type string
		err = rows.Scan(&table_name, &table_type)
		if err != nil {
			log.Fatal(err)
		}
		switch (table_type){
		case "agg_src":
			agg_src_table = table_name
		case "agg_dst":
			agg_dst_table = table_name
		case "agg_srcdst":
			agg_srcdst_table = table_name
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	sql := SQL{
		db_name: db_name,
		metadata_table: metadata_table,
		agg_src_table: agg_src_table,
		agg_dst_table: agg_dst_table,
		agg_srcdst_table: agg_srcdst_table,
		delta_t: delta_t.String(),
	}

	return sql
}