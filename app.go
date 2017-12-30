package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/gorilla/websocket"
	"github.com/jagandecapri/vision/server"
	"encoding/json"
	"github.com/jagandecapri/vision/tree"
)

var upgrader = websocket.Upgrader{} // use default options

type Socket struct{
	URI string
}

type Data struct{
	points []tree.Point
}

func BootServer(data chan server.HttpData) {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveTemplate)

	hub := server.NewHub()
	go hub.Run()
	go func(data chan server.HttpData, hub *server.Hub){
		for{
			select {
			case tmp := <-data:
				json, err := json.Marshal(tmp)
				if err != nil{
					log.Println("JSON encode err: ", err)
				}
				hub.Broadcast <- json
			default:
			}
		}
	}(data, hub)

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(hub, w, r)
	})
	log.Println("Listening...")
	http.ListenAndServe(":3001", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		log.Println("err file: ", err)
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		log.Println("request to dir")
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles(lp, fp))

	data := Socket{URI: "ws://"+r.Host+"/echo"}
	log.Println(data)
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}