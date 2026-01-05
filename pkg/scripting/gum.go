// Package scripting provides Gum-style scriptable TUI commands
package scripting

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Gum provides a Go API for Gum commands
type Gum struct {
	theme string
}

// NewGum creates a new Gum wrapper
func NewGum() *Gum {
	return &Gum{theme: "dracula"}
}

// SetTheme sets the Gum theme
func (g *Gum) SetTheme(theme string) {
	g.theme = theme
}

// Choose presents a selection menu and returns the choice
func (g *Gum) Choose(prompt string, options []string) (string, error) {
	args := []string{"choose"}
	args = append(args, options...)

	return g.run(args...)
}

// Confirm asks for yes/no confirmation
func (g *Gum) Confirm(prompt string) (bool, error) {
	result, err := g.run("confirm", prompt)
	if err != nil {
		// Exit code 1 means "no"
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, err
	}
	return result == "" || strings.TrimSpace(result) == "", nil
}

// Input prompts for text input
func (g *Gum) Input(prompt, placeholder string) (string, error) {
	args := []string{"input", "--placeholder", placeholder}
	if prompt != "" {
		args = append(args, "--header", prompt)
	}
	return g.run(args...)
}

// Write prompts for multi-line text input
func (g *Gum) Write(placeholder string) (string, error) {
	return g.run("write", "--placeholder", placeholder)
}

// Spin shows a spinner while running a command
func (g *Gum) Spin(title string, command string) (string, error) {
	return g.run("spin", "--title", title, "--", "sh", "-c", command)
}

// Filter provides fuzzy filtering of items
func (g *Gum) Filter(items []string, prompt string) (string, error) {
	args := []string{"filter", "--placeholder", prompt}
	
	cmd := exec.Command("gum", args...)
	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))
	
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(stdout.String()), nil
}

// Style applies Gum styling to text
func (g *Gum) Style(text string, opts ...StyleOption) (string, error) {
	args := []string{"style"}
	
	for _, opt := range opts {
		args = append(args, opt.args()...)
	}
	
	args = append(args, text)
	return g.run(args...)
}

// Format formats text with Gum
func (g *Gum) Format(text, formatType string) (string, error) {
	return g.run("format", "-t", formatType, text)
}

// Pager displays content in a pager
func (g *Gum) Pager(content string) error {
	cmd := exec.Command("gum", "pager")
	cmd.Stdin = strings.NewReader(content)
	return cmd.Run()
}

// run executes a gum command
func (g *Gum) run(args ...string) (string, error) {
	cmd := exec.Command("gum", args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("gum %s: %w - %s", args[0], err, stderr.String())
	}
	
	return strings.TrimSpace(stdout.String()), nil
}

// StyleOption defines styling options
type StyleOption struct {
	key   string
	value string
}

func (s StyleOption) args() []string {
	if s.value != "" {
		return []string{"--" + s.key, s.value}
	}
	return []string{"--" + s.key}
}

// Style options
func Bold() StyleOption              { return StyleOption{key: "bold"} }
func Italic() StyleOption            { return StyleOption{key: "italic"} }
func Foreground(c string) StyleOption { return StyleOption{key: "foreground", value: c} }
func Background(c string) StyleOption { return StyleOption{key: "background", value: c} }
func Border(style string) StyleOption { return StyleOption{key: "border", value: style} }
func Padding(p string) StyleOption    { return StyleOption{key: "padding", value: p} }
func Margin(m string) StyleOption     { return StyleOption{key: "margin", value: m} }

// IsGumInstalled checks if gum is available
func IsGumInstalled() bool {
	_, err := exec.LookPath("gum")
	return err == nil
}

