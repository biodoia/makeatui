// Package freeze provides integration with Charm's Freeze for screenshot generation
package freeze

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Config holds freeze configuration
type Config struct {
	Output     string            // output file path (png, svg, webp)
	Theme      string            // syntax theme
	Font       FontConfig        // font settings
	Window     WindowConfig      // window chrome settings
	Padding    [4]int            // padding [top, right, bottom, left]
	Margin     [4]int            // margin
	Border     BorderConfig      // border settings
	Shadow     ShadowConfig      // shadow settings
	Background string            // background color
	Language   string            // syntax highlighting language
	Lines      string            // line range "1-10" or "5,10,15"
	ShowLines  bool              // show line numbers
}

// FontConfig for text rendering
type FontConfig struct {
	Family   string
	Size     float64
	Ligatures bool
}

// WindowConfig for window chrome
type WindowConfig struct {
	Show      bool   // show window decorations
	Title     string // window title
	Style     string // macos, windows, linux
	Controls  bool   // show window controls
}

// BorderConfig for border styling
type BorderConfig struct {
	Radius int
	Width  int
	Color  string
}

// ShadowConfig for drop shadow
type ShadowConfig struct {
	Show   bool
	X      int
	Y      int
	Blur   int
	Color  string
}

// DefaultConfig returns sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Output: "screenshot.png",
		Theme:  "dracula",
		Font: FontConfig{
			Family:   "JetBrains Mono",
			Size:     14,
			Ligatures: true,
		},
		Window: WindowConfig{
			Show:     true,
			Title:    "MakeaTUI",
			Style:    "macos",
			Controls: true,
		},
		Padding:    [4]int{20, 40, 20, 40},
		Margin:     [4]int{0, 0, 0, 0},
		Background: "#1e1e2e",
		Border: BorderConfig{
			Radius: 8,
			Width:  0,
		},
		Shadow: ShadowConfig{
			Show:  true,
			X:     0,
			Y:     8,
			Blur:  16,
			Color: "rgba(0,0,0,0.5)",
		},
		ShowLines: true,
	}
}

// UltravioletConfig returns MakeaTUI themed config
func UltravioletConfig() *Config {
	cfg := DefaultConfig()
	cfg.Theme = "dracula"
	cfg.Background = "#0D0221"
	cfg.Window.Title = "MakeaTUI"
	cfg.Border.Color = "#9D4EDD"
	return cfg
}

// Freezer wraps freeze functionality
type Freezer struct {
	config *Config
	binary string // path to freeze binary
}

// NewFreezer creates a new freezer instance
func NewFreezer(config *Config) *Freezer {
	if config == nil {
		config = DefaultConfig()
	}
	return &Freezer{
		config: config,
		binary: "freeze", // assumes freeze is in PATH
	}
}

// SetBinary sets custom path to freeze binary
func (f *Freezer) SetBinary(path string) *Freezer {
	f.binary = path
	return f
}

// IsAvailable checks if freeze is installed
func (f *Freezer) IsAvailable() bool {
	_, err := exec.LookPath(f.binary)
	return err == nil
}

// CaptureFile captures a file to an image
func (f *Freezer) CaptureFile(inputPath, outputPath string) error {
	args := f.buildArgs(outputPath)
	args = append(args, inputPath)
	return f.run(args...)
}

// CaptureString captures string content to an image
func (f *Freezer) CaptureString(content, outputPath string) error {
	// Write to temp file
	tmpFile, err := os.CreateTemp("", "freeze-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	return f.CaptureFile(tmpFile.Name(), outputPath)
}

// CaptureCommand captures command output to an image
func (f *Freezer) CaptureCommand(command, outputPath string) error {
	args := f.buildArgs(outputPath)
	args = append(args, "--execute", command)
	return f.run(args...)
}

func (f *Freezer) buildArgs(output string) []string {
	cfg := f.config
	args := []string{
		"--output", output,
	}

	if cfg.Theme != "" {
		args = append(args, "--theme", cfg.Theme)
	}
	if cfg.Font.Family != "" {
		args = append(args, "--font.family", cfg.Font.Family)
	}
	if cfg.Font.Size > 0 {
		args = append(args, "--font.size", fmt.Sprintf("%.0f", cfg.Font.Size))
	}
	if cfg.Window.Show {
		args = append(args, "--window")
		if cfg.Window.Title != "" {
			args = append(args, "--window.title", cfg.Window.Title)
		}
	}
	if cfg.Background != "" {
		args = append(args, "--background", cfg.Background)
	}
	if cfg.ShowLines {
		args = append(args, "--show-line-numbers")
	}
	if cfg.Language != "" {
		args = append(args, "--language", cfg.Language)
	}
	if cfg.Lines != "" {
		args = append(args, "--lines", cfg.Lines)
	}
	if cfg.Border.Radius > 0 {
		args = append(args, "--border.radius", fmt.Sprintf("%d", cfg.Border.Radius))
	}
	if cfg.Shadow.Show {
		args = append(args, "--shadow.blur", fmt.Sprintf("%d", cfg.Shadow.Blur))
	}

	// Padding
	padding := fmt.Sprintf("%d %d %d %d", 
		cfg.Padding[0], cfg.Padding[1], cfg.Padding[2], cfg.Padding[3])
	args = append(args, "--padding", strings.TrimSpace(padding))

	return args
}

func (f *Freezer) run(args ...string) error {
	cmd := exec.Command(f.binary, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("freeze failed: %s: %w", stderr.String(), err)
	}
	return nil
}

