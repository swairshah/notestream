package transcribe

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/swairshah/notestream/internal/config"
)

type Transcriber struct {
	config *config.Config
}

func NewTranscriber(cfg *config.Config) *Transcriber {
	return &Transcriber{
		config: cfg,
	}
}

func (t *Transcriber) TranscribeFile(inputFile string) (string, error) {
	consolidatedFile := filepath.Join(t.config.OutputDir, "consolidated_transcriptions.txt")

	cmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("%s -f %s -np", t.config.WhisperPath, inputFile),
	)

	fmt.Println("cmd:", cmd)

	// Open the consolidated file in append mode
	outFile, err := os.OpenFile(consolidatedFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open consolidated file: %w", err)
	}
	defer outFile.Close()

	// Write a header for this transcription
	_, err = fmt.Fprintf(outFile, "\n--- Transcription for %s ---\n", filepath.Base(inputFile))
	if err != nil {
		return "", fmt.Errorf("failed to write transcription header: %w", err)
	}

	// Set up command to pipe stdout to the file
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr // Redirect stderr to os.Stderr for error logging

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	return consolidatedFile, nil
}
