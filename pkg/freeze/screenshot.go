// Package freeze - Screenshot utilities for TUI
package freeze

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ScreenshotFormat defines output format
type ScreenshotFormat string

const (
	FormatPNG  ScreenshotFormat = "png"
	FormatSVG  ScreenshotFormat = "svg"
	FormatWEBP ScreenshotFormat = "webp"
	FormatJPG  ScreenshotFormat = "jpg"
)

// Screenshot captures a screenshot with automatic naming
type Screenshot struct {
	freezer   *Freezer
	outputDir string
	prefix    string
	format    ScreenshotFormat
}

// NewScreenshot creates a screenshot helper
func NewScreenshot(config *Config) *Screenshot {
	return &Screenshot{
		freezer:   NewFreezer(config),
		outputDir: "screenshots",
		prefix:    "makeatui",
		format:    FormatPNG,
	}
}

// SetOutputDir sets the output directory
func (s *Screenshot) SetOutputDir(dir string) *Screenshot {
	s.outputDir = dir
	return s
}

// SetPrefix sets the filename prefix
func (s *Screenshot) SetPrefix(prefix string) *Screenshot {
	s.prefix = prefix
	return s
}

// SetFormat sets the output format
func (s *Screenshot) SetFormat(format ScreenshotFormat) *Screenshot {
	s.format = format
	return s
}

// generateFilename creates a unique filename
func (s *Screenshot) generateFilename(name string) string {
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s_%s_%s.%s", s.prefix, name, timestamp, s.format)
	return filepath.Join(s.outputDir, filename)
}

// ensureDir creates output directory if needed
func (s *Screenshot) ensureDir() error {
	return os.MkdirAll(s.outputDir, 0755)
}

// CaptureView captures a TUI view string
func (s *Screenshot) CaptureView(view, name string) (string, error) {
	if err := s.ensureDir(); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	outputPath := s.generateFilename(name)
	if err := s.freezer.CaptureString(view, outputPath); err != nil {
		return "", err
	}
	return outputPath, nil
}

// CaptureCode captures code/source file
func (s *Screenshot) CaptureCode(code, language, name string) (string, error) {
	if err := s.ensureDir(); err != nil {
		return "", err
	}

	// Create temp file with proper extension
	ext := language
	if ext == "" {
		ext = "txt"
	}
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("code-*.%s", ext))
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(code); err != nil {
		return "", err
	}
	tmpFile.Close()

	// Configure language highlighting
	s.freezer.config.Language = language

	outputPath := s.generateFilename(name)
	if err := s.freezer.CaptureFile(tmpFile.Name(), outputPath); err != nil {
		return "", err
	}
	return outputPath, nil
}

// CaptureProgram runs and captures a program's output
func (s *Screenshot) CaptureProgram(command, name string) (string, error) {
	if err := s.ensureDir(); err != nil {
		return "", err
	}

	outputPath := s.generateFilename(name)
	if err := s.freezer.CaptureCommand(command, outputPath); err != nil {
		return "", err
	}
	return outputPath, nil
}

// Batch captures multiple screenshots
type Batch struct {
	screenshot *Screenshot
	items      []BatchItem
}

// BatchItem represents a single capture item
type BatchItem struct {
	Type    string // "view", "code", "command"
	Content string
	Name    string
	Lang    string // for code
}

// NewBatch creates a batch capture
func NewBatch(config *Config) *Batch {
	return &Batch{
		screenshot: NewScreenshot(config),
		items:      []BatchItem{},
	}
}

// AddView adds a view to capture
func (b *Batch) AddView(view, name string) *Batch {
	b.items = append(b.items, BatchItem{Type: "view", Content: view, Name: name})
	return b
}

// AddCode adds code to capture
func (b *Batch) AddCode(code, language, name string) *Batch {
	b.items = append(b.items, BatchItem{Type: "code", Content: code, Name: name, Lang: language})
	return b
}

// AddCommand adds a command to capture
func (b *Batch) AddCommand(command, name string) *Batch {
	b.items = append(b.items, BatchItem{Type: "command", Content: command, Name: name})
	return b
}

// Execute runs all captures and returns paths
func (b *Batch) Execute() ([]string, error) {
	var paths []string
	for _, item := range b.items {
		var path string
		var err error

		switch item.Type {
		case "view":
			path, err = b.screenshot.CaptureView(item.Content, item.Name)
		case "code":
			path, err = b.screenshot.CaptureCode(item.Content, item.Lang, item.Name)
		case "command":
			path, err = b.screenshot.CaptureProgram(item.Content, item.Name)
		}

		if err != nil {
			return paths, fmt.Errorf("failed to capture %s: %w", item.Name, err)
		}
		paths = append(paths, path)
	}
	return paths, nil
}

