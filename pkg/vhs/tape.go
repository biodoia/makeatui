// Package vhs provides VHS tape generation for recording TUI sessions
package vhs

import (
	"time"
)

// Tape represents a VHS recording tape
type Tape struct {
	Output      string
	Settings    Settings
	Instructions []Instruction
}

// Settings defines VHS recording settings
type Settings struct {
	Width         int
	Height        int
	FontFamily    string
	FontSize      int
	Theme         string
	Padding       int
	Framerate     int
	PlaybackSpeed float64
	TypingSpeed   time.Duration
}

// DefaultSettings returns sensible defaults
func DefaultSettings() Settings {
	return Settings{
		Width:         120,
		Height:        40,
		FontFamily:    "JetBrains Mono",
		FontSize:      14,
		Theme:         "Dracula",
		Padding:       20,
		Framerate:     60,
		PlaybackSpeed: 1.0,
		TypingSpeed:   50 * time.Millisecond,
	}
}

// MakeaTUISettings returns settings optimized for MakeaTUI demos
func MakeaTUISettings() Settings {
	s := DefaultSettings()
	s.Theme = "Catppuccin Mocha"
	s.Width = 140
	s.Height = 45
	return s
}

// Instruction represents a VHS instruction
type Instruction struct {
	Type    InstructionType
	Content string
	Delay   time.Duration
}

// InstructionType defines VHS instruction types
type InstructionType string

const (
	InstType      InstructionType = "Type"
	InstEnter     InstructionType = "Enter"
	InstSleep     InstructionType = "Sleep"
	InstCtrl      InstructionType = "Ctrl"
	InstDown      InstructionType = "Down"
	InstUp        InstructionType = "Up"
	InstLeft      InstructionType = "Left"
	InstRight     InstructionType = "Right"
	InstTab       InstructionType = "Tab"
	InstEscape    InstructionType = "Escape"
	InstBackspace InstructionType = "Backspace"
	InstHide      InstructionType = "Hide"
	InstShow      InstructionType = "Show"
	InstSetTypingSpeed InstructionType = "Set TypingSpeed"
)

// NewTape creates a new VHS tape
func NewTape(output string) *Tape {
	return &Tape{
		Output:       output,
		Settings:     DefaultSettings(),
		Instructions: []Instruction{},
	}
}

// SetSettings updates the tape settings
func (t *Tape) SetSettings(s Settings) *Tape {
	t.Settings = s
	return t
}

// Type adds a type instruction
func (t *Tape) Type(text string) *Tape {
	t.Instructions = append(t.Instructions, Instruction{
		Type:    InstType,
		Content: text,
	})
	return t
}

// Enter adds an enter keystroke
func (t *Tape) Enter() *Tape {
	t.Instructions = append(t.Instructions, Instruction{Type: InstEnter})
	return t
}

// Sleep adds a pause
func (t *Tape) Sleep(d time.Duration) *Tape {
	t.Instructions = append(t.Instructions, Instruction{
		Type:    InstSleep,
		Content: d.String(),
	})
	return t
}

// Ctrl adds a control key combination
func (t *Tape) Ctrl(key string) *Tape {
	t.Instructions = append(t.Instructions, Instruction{
		Type:    InstCtrl,
		Content: key,
	})
	return t
}

// Down adds a down arrow press
func (t *Tape) Down(n int) *Tape {
	for i := 0; i < n; i++ {
		t.Instructions = append(t.Instructions, Instruction{Type: InstDown})
	}
	return t
}

// Up adds an up arrow press
func (t *Tape) Up(n int) *Tape {
	for i := 0; i < n; i++ {
		t.Instructions = append(t.Instructions, Instruction{Type: InstUp})
	}
	return t
}

// Tab adds a tab press
func (t *Tape) Tab() *Tape {
	t.Instructions = append(t.Instructions, Instruction{Type: InstTab})
	return t
}

// Hide hides cursor/output temporarily
func (t *Tape) Hide() *Tape {
	t.Instructions = append(t.Instructions, Instruction{Type: InstHide})
	return t
}

// Show shows cursor/output
func (t *Tape) Show() *Tape {
	t.Instructions = append(t.Instructions, Instruction{Type: InstShow})
	return t
}

// SetTypingSpeed changes typing speed
func (t *Tape) SetTypingSpeed(d time.Duration) *Tape {
	t.Instructions = append(t.Instructions, Instruction{
		Type:    InstSetTypingSpeed,
		Content: d.String(),
	})
	return t
}

