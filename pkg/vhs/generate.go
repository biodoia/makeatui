// Package vhs - Tape generation and execution
package vhs

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

const tapeTemplate = `# MakeaTUI VHS Tape
# Generated automatically

Output {{ .Output }}

Set Width {{ .Settings.Width }}
Set Height {{ .Settings.Height }}
Set FontFamily "{{ .Settings.FontFamily }}"
Set FontSize {{ .Settings.FontSize }}
Set Theme "{{ .Settings.Theme }}"
Set Padding {{ .Settings.Padding }}
Set Framerate {{ .Settings.Framerate }}
Set PlaybackSpeed {{ .Settings.PlaybackSpeed }}

{{ range .Instructions }}
{{- if eq .Type "Type" }}Type "{{ .Content }}"
{{ else if eq .Type "Sleep" }}Sleep {{ .Content }}
{{ else if eq .Type "Ctrl" }}Ctrl+{{ .Content }}
{{ else if eq .Type "Set TypingSpeed" }}Set TypingSpeed {{ .Content }}
{{ else }}{{ .Type }}
{{ end }}
{{- end }}`

// Generate creates the VHS tape file content
func (t *Tape) Generate() (string, error) {
	tmpl, err := template.New("tape").Parse(tapeTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, t); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Save writes the tape to a file
func (t *Tape) Save(path string) error {
	content, err := t.Generate()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// Run executes the tape with VHS
func (t *Tape) Run() error {
	// Save to temp file
	tmpFile, err := os.CreateTemp("", "makeatui-*.tape")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if err := t.Save(tmpFile.Name()); err != nil {
		return err
	}

	// Run VHS
	cmd := exec.Command("vhs", tmpFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// IsVHSInstalled checks if VHS is available
func IsVHSInstalled() bool {
	_, err := exec.LookPath("vhs")
	return err == nil
}

// CreateDemoTape creates a tape demonstrating MakeaTUI
func CreateDemoTape(outputPath string) *Tape {
	return NewTape(outputPath).
		SetSettings(MakeaTUISettings()).
		Hide().
		Type("./makeatui").
		Enter().
		Sleep(2 * time.Second).
		Show().
		Sleep(1 * time.Second).
		// Navigate sidebar
		Down(2).
		Sleep(500 * time.Millisecond).
		Enter().
		Sleep(1 * time.Second).
		// Switch to canvas
		Tab().
		Sleep(500 * time.Millisecond).
		// Move around
		Down(5).
		Sleep(300 * time.Millisecond).
		Down(5).
		Sleep(300 * time.Millisecond).
		// Add another component
		Tab().
		Sleep(500 * time.Millisecond).
		Down(1).
		Sleep(300 * time.Millisecond).
		Enter().
		Sleep(1 * time.Second).
		// Show help
		Type("?").
		Sleep(2 * time.Second).
		Type("?").
		Sleep(500 * time.Millisecond).
		// Quit
		Type("q").
		Sleep(1 * time.Second)
}

// CreateTutorialTape creates a tutorial tape
func CreateTutorialTape(outputPath string) *Tape {
	return NewTape(outputPath).
		SetSettings(MakeaTUISettings()).
		SetTypingSpeed(75 * time.Millisecond).
		Type("# Welcome to MakeaTUI!").
		Enter().
		Sleep(1 * time.Second).
		Type("# Let's create a beautiful TUI").
		Enter().
		Sleep(1 * time.Second).
		Type("./makeatui").
		Enter().
		Sleep(2 * time.Second).
		// More tutorial steps...
		Type("q").
		Sleep(500 * time.Millisecond)
}

// GitHubActionConfig generates VHS GitHub Action configuration
func GitHubActionConfig() string {
	return `name: Generate Demo GIF

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  vhs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Build MakeaTUI
        run: go build -o makeatui .

      - name: Generate Demo GIF
        uses: charmbracelet/vhs-action@v2
        with:
          path: 'demo.tape'

      - name: Upload GIF
        uses: actions/upload-artifact@v4
        with:
          name: demo
          path: demo.gif
`
}

// ParseTape parses an existing tape file
func ParseTape(content string) (*Tape, error) {
	tape := NewTape("")
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "Output":
			if len(parts) > 1 {
				tape.Output = parts[1]
			}
		case "Type":
			if len(parts) > 1 {
				// Remove quotes
				text := strings.Trim(parts[1], "\"")
				tape.Type(text)
			}
		case "Enter":
			tape.Enter()
		case "Tab":
			tape.Tab()
		case "Down":
			tape.Down(1)
		case "Up":
			tape.Up(1)
		case "Sleep":
			if len(parts) > 1 {
				tape.Instructions = append(tape.Instructions, Instruction{
					Type:    InstSleep,
					Content: parts[1],
				})
			}
		}
	}

	return tape, nil
}

