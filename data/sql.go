package data

import (
	"log"
	"database/sql"
	"time"
	"html/template"
	"bytes"
	"github.com/jagandecapri/vision/preprocess"
	"github.com/jagandecapri/vision/tree"
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
	metadata_table := `CREATE TABLE IF NOT EXISTS` +s.metadata_table+ ` (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
		delta_t	TEXT,
		table_name TEXT,
		time NUMERIC)`

	_, err = db.Exec(metadata_table)
	if err != nil {
		log.Printf("%q: %s\n", err, metadata_table)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO metadata_table(delta_t, table_name, time) values(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(s.delta_t, s.agg_src_table, s.time); err != nil {
		log.Fatal(err)
	}
	if _, err := stmt.Exec(s.delta_t, s.agg_dst_table, s.time); err != nil {
		log.Fatal(err)
	}
	if _, err := stmt.Exec(s.delta_t, s.agg_srcdst_table, s.time); err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	//CREATE TABLES TO STORE DATA
	t := template.New("create tables")
	t, _ = t.Parse(`CREATE TABLE IF NOT EXISTS {{.TableName}} (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
		batch INTEGER,
		flow_key TEXT,
		{{range $column_name, $type := .FlowKeys}}{{$column_name}} {{$type}},
		{{end}}
		last_packet_timestamp	NUMERIC
	)`)

	var tpl bytes.Buffer
	var result string

	if err := t.Execute(&tpl, struct{TableName string
	FlowKeys map[string]string}{
		TableName: s.agg_src_table,
		FlowKeys: s.flow_keys}); err != nil {
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

func (s *SQL) WriteToDb(acc_c preprocess.AccumulatorChannels, done chan struct{}){
	db, err := sql.Open("sqlite3", s.db_name)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//CREATE TABLES TO STORE DATA
	t := template.New("insert data to tables").Funcs(template.FuncMap{"getKeys": func(flow_keys map[string]string) string{
			tmp := []string{}
			for column_name, _ := range flow_keys{
				tmp = append(tmp, column_name)
			}
			ret := strings.Join(tmp, ",")
			return ret
		},
	})
	t, _ = t.Parse(`INSERT INTO {{.TableName}}(getKeys .FlowKeys) VALUES ({{range $column_name, $type := .FlowKeys}},{{end}})`)


	var tpl bytes.Buffer
	var result string

	if err := t.Execute(&tpl, struct{FlowKeys map[string]string}{FlowKeys: s.flow_keys}); err != nil {
		log.Fatal(err)
	}

	result = tpl.String()

	go func(){
		var tmp []tree.Point

		for{
			select{
				case X := <-acc_c.AggSrc:
					tmp = X
				case X := <-acc_c.AggDst:
					tmp = X
				case X := <-acc_c.AggSrcDst:
					tmp = X
				case <-done:
					return
				default:
			}

			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}

			stmt, err := tx.Prepare(result)
			if err != nil {
				log.Fatal(err)
			}

			for _, p := range tmp {
				var tmp1 = make([]interface{}, len(p.Vec_map))
				for _, val := range p.Vec_map{
					tmp1 = append(tmp1, val)
				}

				_, err = stmt.Exec(tmp1...)
				if err != nil {
					log.Fatal(err)
				}
			}
			tx.Commit()
			stmt.Close()
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
		flow_keys: map[string]string{"nbPacket": "REAL",
		"nbSrcPort": "REAL",
		"nbDstPort": "REAL",
		"nbSrcs": "REAL",
		"nbDsts": "REAL",
		"perSYN": "REAL",
		"perACK": "REAL",
		"perICMP": "REAL",
		"perRST": "REAL",
		"perFIN": "REAL",
		"perCWR": "REAL",
		"perURG": "REAL",
		"avgPktSize": "REAL",
		"meanTTL": "REAL"},
		delta_t: delta_t.String(),
		time: now_string,
	}

	sql.SetupDb()

	return sql
}
