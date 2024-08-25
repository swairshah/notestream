package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/swairshah/notestream/internal/audio"
	"github.com/swairshah/notestream/internal/config"
	"github.com/swairshah/notestream/internal/transcribe"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	transcriber := transcribe.NewTranscriber(cfg)
	capture, err := audio.NewCapture(cfg, transcriber)

	if err != nil {
		fmt.Println("Error initializing audio capture:", err)
		return
	}
	defer capture.Close()

	// Start capturing audio
	err = capture.Start()
	if err != nil {
		fmt.Println("Error starting audio capture:", err)
		return
	}

	// Wait for interrupt signal to stop
	fmt.Println("Recording... Press Ctrl+C to stop.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Stop capturing audio
	err = capture.Stop()
	if err != nil {
		fmt.Println("Error stopping audio capture:", err)
	}

	fmt.Println("Recording stopped.")

}
