package main

import (
	"testing"
	"github.com/jagandecapri/vision/server"
	"os"
	"os/signal"
)

func TestBootServer(t *testing.T) {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.

	test_data := server.HttpData{
		server.Graph{
			Graph_metadata: server.Graph_metadata{ID: "first-second",
				Column_x: "first",
				Column_y: "second"},
			PointsContainer: []server.PointsContainer{{
				Point_list: []server.Point{{
					Point_data: server.Point_data{X: 0.05, Y: 0.05},
				}},
				Points_metadata: server.Points_metadata{Color: "#FF0000"},
			},
				{
					Point_list: []server.Point{{
						Point_data: server.Point_data{X: 0.15, Y: 0.15},
					}},
					Points_metadata: server.Points_metadata{Color: "#800000"},
				},
				{
					Point_list: []server.Point{{
						Point_data: server.Point_data{X: 0.25, Y: 0.25},
					}},
					Points_metadata: server.Points_metadata{Color: "#FFFF00"},
				},
			},
		},
		server.Graph{
			Graph_metadata: server.Graph_metadata{ID: "third-fourth",
				Column_x: "third",
				Column_y: "fourth"},
			PointsContainer: []server.PointsContainer{{
				Point_list: []server.Point{{
					Point_data: server.Point_data{X: 0.05, Y: 0.05},
				}},
				Points_metadata: server.Points_metadata{Color: "#FF0000"},
			},
				{
					Point_list: []server.Point{{
						Point_data: server.Point_data{X: 0.15, Y: 0.15},
					}},
					Points_metadata: server.Points_metadata{Color: "#800000"},
				},
				{
					Point_list: []server.Point{{
						Point_data: server.Point_data{X: 0.25, Y: 0.25},
					}},
					Points_metadata: server.Points_metadata{Color: "#FFFF00"},
				},
			},
		},
	}

	data := make(chan server.HttpData)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	BootServer(data)

	//ticker := time.NewTicker(10 * time.Second)
	//
	//go func() {
	//	for t := range ticker.C{
	//		fmt.Println("Tick at", t)
	//		data <- test_data
	//	}
	//}()
	//
	//for {
	//	// Block until a signal is received.
	//	select {
	//		case <- c:
	//			ticker.Stop()
	//			break;
	//		default:
	//	}
	//}

	for {
		data <- test_data
		select {
			case <-c:
				break;
			default:
		}
	}
}
