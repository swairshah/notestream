package audio

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
	"github.com/swairshah/notestream/internal/config"
	"github.com/swairshah/notestream/internal/transcribe"
)

type Capture struct {
	config      *config.Config
	stream      *portaudio.Stream
	buffer      []int16
	chunkIndex  int
	transcriber *transcribe.Transcriber
	wg          sync.WaitGroup
}

func NewCapture(cfg *config.Config, transcriber *transcribe.Transcriber) (*Capture, error) {
	err := portaudio.Initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PortAudio: %w", err)
	}

	bufferSize := cfg.SampleRate * cfg.Channels * int(cfg.RecordInterval.Seconds())
	return &Capture{
		config:      cfg,
		buffer:      make([]int16, bufferSize),
		transcriber: transcriber,
	}, nil
}

func (c *Capture) Start() error {
	var err error
	c.stream, err = portaudio.OpenDefaultStream(c.config.Channels, 0, float64(c.config.SampleRate), len(c.buffer), c.processAudio)
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %w", err)
	}

	err = c.stream.Start()
	if err != nil {
		return fmt.Errorf("failed to start audio stream: %w", err)
	}

	go c.saveChunks()

	return nil
}

func (c *Capture) processAudio(in []int16) {
	copy(c.buffer, in)
}

func (c *Capture) saveChunks() {
	ticker := time.NewTicker(c.config.RecordInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.saveChunk()
		c.chunkIndex++
	}
}

func (c *Capture) saveChunk() {
	filename := filepath.Join(c.config.OutputDir, fmt.Sprintf("chunk_%d.wav", c.chunkIndex))

	err := os.MkdirAll(c.config.OutputDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	// Create a new encoder
	enc := wav.NewEncoder(file, c.config.SampleRate, 16, c.config.Channels, 1)
	defer enc.Close()

	// Convert int16 samples to int
	audioData := make([]int, len(c.buffer))
	for i, sample := range c.buffer {
		audioData[i] = int(sample)
	}

	// Create audio.IntBuffer
	buf := &audio.IntBuffer{
		Data: audioData,
		Format: &audio.Format{
			NumChannels: c.config.Channels,
			SampleRate:  c.config.SampleRate,
		},
	}

	// Write audio data
	if err := enc.Write(buf); err != nil {
		fmt.Printf("Failed to write audio data: %v\n", err)
		return
	}

	fmt.Printf("Saved audio chunk: %s\n", filename)

	c.wg.Add(1)
	go func(inputFile string) {
		defer c.wg.Done()
		outputFile, err := c.transcriber.TranscribeFile(inputFile)
		if err != nil {
			fmt.Printf("Failed to transcribe %s: %v\n", inputFile, err)
		}
		fmt.Printf("Transcription saved to %s\n", outputFile)
	}(filename)
}

func (c *Capture) Stop() error {
	if c.stream != nil {
		err := c.stream.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop audio stream: %w", err)
		}
	}
	c.wg.Wait()

	return nil
}

func (c *Capture) Close() error {
	if c.stream != nil {
		err := c.stream.Close()
		if err != nil {
			return fmt.Errorf("failed to close audio stream: %w", err)
		}
	}
	portaudio.Terminate()
	return nil
}
