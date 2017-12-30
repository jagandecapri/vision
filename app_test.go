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
	data := make(chan server.HttpData)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go BootServer(data)

	// Block until a signal is received.
	<-c
}
