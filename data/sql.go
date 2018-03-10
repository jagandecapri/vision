package data

import (
	"log"
	"database/sql"
	"time"
	"html/template"
	"bytes"
	"github.com/jagandecapri/vision/preprocess"
)

type SQL struct{
	db_name string
	metadata_table string
	agg_src_table string
	agg_dst_table string
	agg_srcdst_table string
}

func (s *SQL) SetupDb(delta_t time.Duration){
	db, err := sql.Open("sqlite3", s.db_name)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//CREATE TABLE TO STORE METADATA
	metadata_table := `CREATE TABLE IF NOT EXISTS` +s.metadata_table+ ` (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
		delta_t	TEXT)`

	_, err = db.Exec(metadata_table)
	if err != nil {
		log.Printf("%q: %s\n", err, metadata_table)
		return
	}

	_, err = db.Exec("insert into metadata_table(delta_t) values('" +delta_t.String()+ "')")
	if err != nil {
		log.Fatal(err)
	}

	//CREATE TABLES TO STORE DATA
	t := template.New("create tables")
	t, _ = t.Parse(`CREATE TABLE IF NOT EXISTS {{.TableName}} (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
		flow_data	TEXT,
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
		log.Printf("%q: %s\n", err, result)
		return
	}

	if err := t.Execute(&tpl, struct{TableName string}{TableName: s.agg_dst_table}); err != nil {
		log.Println(err)
	}

	result = tpl.String()

	_, err = db.Exec(result)
	if err != nil {
		log.Printf("%q: %s\n", err, result)
		return
	}

	if err := t.Execute(&tpl, struct{TableName string}{TableName: s.agg_srcdst_table}); err != nil {
		log.Println(err)
	}

	result = tpl.String()

	_, err = db.Exec(result)
	if err != nil {
		log.Printf("%q: %s\n", err, result)
		return
	}
}

func WriteToDb(acc_c preprocess.AccumulatorChannel, done chan struct{}){
	go func(){
		for{
			select{
				case X := <-acc_c:

				case <-done:
					return
				default:
			}
		}
	}()
}

func NewSQL(acc_c preprocess.AccumulatorChannels, delta_t time.Duration) SQL{
	now := time.Now()
	now_string := now.Format(time.RFC3339)
	sql := SQL{
		db_name: "vision.db",
		metadata_table: "metadata_" + now_string,
		agg_src_table: "agg_src_" + now_string,
		agg_dst_table: "agg_dst_" + now_string,
		agg_srcdst_table: "agg_srcdst_" + now_string,
	}

	sql.SetupDb(delta_t)

	return sql
}
