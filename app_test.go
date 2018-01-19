package main

import (
	"testing"
	"github.com/jagandecapri/vision/server"
	"os"
	"os/signal"
	"github.com/jagandecapri/vision/tree"
	"encoding/json"
	"log"
)

func TestMarshalData(t *testing.T){
	point := tree.Point{Id: 1, Vec_map: map[string]float64{
		"col_1": 0.12345,
		"col_2": 0.56789,

	}}

	point_container := tree.PointContainer{Unit_id: 5,
		Vec: []float64{0.12345, 0.56789},
		Point: point}

	center_point_container := tree.PointContainer{Unit_id: 5, Vec: []float64{0.12345, 0.56789}}

	range_1 := tree.Range{Low: [2]float64{0.1, 0.5},
		High: [2]float64{0.2, 0.6}}

	unit_1 := tree.Unit{Id: 5,
		Cluster_id: 3,
		Dimension: 2,
		Center: center_point_container,
		Points: map[int]tree.PointContainer{1: point_container},
		Center_calculated: true,
		Range: range_1,
	}

	//grid := tree.Grid{Store: map[tree.Range]*tree.Unit{range_1: &unit_1}}

	res, err := json.Marshal(unit_1)
	log.Println(res, err)

}
func TestBootServer(t *testing.T) {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.

	data := make(chan server.HttpData)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go BootServer(data)

	// Block until a signal is received.
	<-c
}
