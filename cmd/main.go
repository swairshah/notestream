package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	hook "github.com/robotn/gohook"
	"github.com/swairshah/notestream/internal/audio"
	"github.com/swairshah/notestream/internal/config"
	"github.com/swairshah/notestream/internal/transcribe"
)

var (
	capture     *audio.Capture
	isRecording bool
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	transcriber := transcribe.NewTranscriber(cfg)
	capture, err = audio.NewCapture(cfg, transcriber)

	if err != nil {
		fmt.Println("Error initializing audio capture:", err)
		return
	}
	defer capture.Close()

	setupHotkey()

	fmt.Println("NoteStream is running. Press Cmd+Shift+R to start/stop recording. Ctrl+C to exit.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("Exiting NoteStream.")

	/*
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
	*/

}

func setupHotkey() {
	go func() {
		fmt.Println("Setting up hotkey Cmd+Shift+R to toggle recording")
		hook.Register(hook.KeyDown, []string{"r", "cmd", "shift"}, func(e hook.Event) {
			toggleRecording()
		})

		s := hook.Start()
		<-hook.Process(s)
	}()
}

func toggleRecording() {
	if capture == nil {
		fmt.Println("Error: Audio capture not initialized")
		return
	}

	if isRecording {
		stopRecording()
	} else {
		startRecording()
	}
}

func startRecording() {
	if err := capture.Start(); err != nil {
		fmt.Println("Error starting audio capture:", err)
		return
	}
	isRecording = true
	fmt.Println("Recording started.")
}

func stopRecording() {
	if err := capture.Stop(); err != nil {
		fmt.Println("Error stopping audio capture:", err)
		return
	}
	isRecording = false
	fmt.Println("Recording stopped.")
}
