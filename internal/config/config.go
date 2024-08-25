package config

import "time"

type Config struct {
	SampleRate     int
	Channels       int
	RecordInterval time.Duration
	OutputDir      string
	WhisperPath    string
}

func Load() (*Config, error) {
	return &Config{
		SampleRate:     44100,
		Channels:       1,
		RecordInterval: 10 * time.Second,
		OutputDir:      "./output",
		WhisperPath:    "./whisper-tiny.en.llamafile",
	}, nil
}
