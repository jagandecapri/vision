package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/gorilla/websocket"
	"fmt"
	"github.com/jagandecapri/vision/orunada/grid"
	"github.com/jagandecapri/vision/orunada/utils"
)

var upgrader = websocket.Upgrader{} // use default options

type Socket struct{
	URI string
}

type Data struct{
	points []grid.Point
}

func BootServer(data chan grid.HttpData) {
	bcast := &utils.ThreadSafeSlice{
		Workers: []*utils.Worker{},
	}

	go utils.Broadcaster(bcast, data)
	quit := make(chan bool)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request){
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer conn.Close()

		wk := &utils.Worker{
			Source: make(chan grid.HttpData),
		}
		bcast.Push(wk)
		q := false
		for !q {
			select {
			case v := <-wk.Source:
				err := conn.WriteJSON(v)
				if err != nil{
					log.Println(err)
				}
			case q = <-quit:
				break
			}
		}
	})
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles(lp, fp))

	data := Socket{URI: "ws://"+r.Host+"/echo"}
	fmt.Println(data)
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}