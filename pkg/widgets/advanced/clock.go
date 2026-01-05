// Package advanced - Clock widget with big digits
package advanced

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Clock provides a digital clock widget
type Clock struct {
	ID         string
	Time       time.Time
	ShowDate   bool
	Show24Hour bool
	ShowSecs   bool
	BigDigits  bool
	style      ClockStyle
}

// ClockStyle holds styling
type ClockStyle struct {
	Container lipgloss.Style
	Digit     lipgloss.Style
	Separator lipgloss.Style
	Date      lipgloss.Style
	AMPM      lipgloss.Style
}

// DefaultClockStyle returns default styling
func DefaultClockStyle() ClockStyle {
	return ClockStyle{
		Container: lipgloss.NewStyle().
			Padding(1, 2),
		Digit: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")),
		Date: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")),
		AMPM: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")).
			Italic(true),
	}
}

// BigDigits font (5 rows high)
var BigDigitFont = map[rune][]string{
	'0': {"█▀▀█", "█  █", "█  █", "█  █", "█▄▄█"},
	'1': {"  ▄█", "   █", "   █", "   █", "   █"},
	'2': {"█▀▀█", "   █", "█▀▀█", "█   ", "█▄▄█"},
	'3': {"█▀▀█", "   █", "█▀▀█", "   █", "█▄▄█"},
	'4': {"█  █", "█  █", "█▀▀█", "   █", "   █"},
	'5': {"█▀▀█", "█   ", "█▀▀█", "   █", "█▄▄█"},
	'6': {"█▀▀█", "█   ", "█▀▀█", "█  █", "█▄▄█"},
	'7': {"█▀▀█", "   █", "   █", "   █", "   █"},
	'8': {"█▀▀█", "█  █", "█▀▀█", "█  █", "█▄▄█"},
	'9': {"█▀▀█", "█  █", "█▀▀█", "   █", "█▄▄█"},
	':': {"    ", "  ▄ ", "    ", "  ▄ ", "    "},
}

// NewClock creates a clock widget
func NewClock(id string) *Clock {
	return &Clock{
		ID:         id,
		Time:       time.Now(),
		ShowDate:   true,
		Show24Hour: true,
		ShowSecs:   true,
		BigDigits:  false,
		style:      DefaultClockStyle(),
	}
}

// SetShow24Hour sets 24-hour mode
func (c *Clock) SetShow24Hour(show24 bool) *Clock {
	c.Show24Hour = show24
	return c
}

// SetShowSecs sets seconds display
func (c *Clock) SetShowSecs(show bool) *Clock {
	c.ShowSecs = show
	return c
}

// SetBigDigits sets big digit mode
func (c *Clock) SetBigDigits(big bool) *Clock {
	c.BigDigits = big
	return c
}

// Update handles tick messages
func (c *Clock) Update(msg tea.Msg) (*Clock, tea.Cmd) {
	switch msg.(type) {
	case ClockTickMsg:
		c.Time = time.Now()
	}
	return c, nil
}

// TickCmd returns animation command
func (c *Clock) TickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return ClockTickMsg{}
	})
}

// ClockTickMsg is sent every second
type ClockTickMsg struct{}

// View renders the clock
func (c *Clock) View() string {
	var parts []string

	// Format time string
	hour := c.Time.Hour()
	ampm := ""
	if !c.Show24Hour {
		if hour >= 12 {
			ampm = " PM"
			if hour > 12 {
				hour -= 12
			}
		} else {
			ampm = " AM"
			if hour == 0 {
				hour = 12
			}
		}
	}

	timeStr := fmt.Sprintf("%02d:%02d", hour, c.Time.Minute())
	if c.ShowSecs {
		timeStr += fmt.Sprintf(":%02d", c.Time.Second())
	}

	if c.BigDigits {
		// Render big digits
		var lines [5][]string
		for _, ch := range timeStr {
			if font, ok := BigDigitFont[ch]; ok {
				for i, row := range font {
					lines[i] = append(lines[i], c.style.Digit.Render(row))
				}
			}
		}

		var bigTime []string
		for _, row := range lines {
			bigTime = append(bigTime, strings.Join(row, " "))
		}
		parts = append(parts, strings.Join(bigTime, "\n"))
	} else {
		// Simple time display
		timeView := c.style.Digit.Render(timeStr)
		if ampm != "" {
			timeView += c.style.AMPM.Render(ampm)
		}
		parts = append(parts, timeView)
	}

	// Date
	if c.ShowDate {
		dateStr := c.Time.Format("Monday, January 2, 2006")
		parts = append(parts, c.style.Date.Render(dateStr))
	}

	return c.style.Container.Render(lipgloss.JoinVertical(lipgloss.Center, parts...))
}

// Stopwatch provides a stopwatch widget
type Stopwatch struct {
	ID        string
	Running   bool
	Elapsed   time.Duration
	Laps      []time.Duration
	startTime time.Time
	style     lipgloss.Style
}

// NewStopwatch creates a stopwatch
func NewStopwatch(id string) *Stopwatch {
	return &Stopwatch{
		ID:   id,
		Laps: []time.Duration{},
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true),
	}
}

// Start starts the stopwatch
func (s *Stopwatch) Start() *Stopwatch {
	if !s.Running {
		s.startTime = time.Now()
		s.Running = true
	}
	return s
}

// Stop stops the stopwatch
func (s *Stopwatch) Stop() *Stopwatch {
	if s.Running {
		s.Elapsed += time.Since(s.startTime)
		s.Running = false
	}
	return s
}

// Toggle toggles running state
func (s *Stopwatch) Toggle() *Stopwatch {
	if s.Running {
		return s.Stop()
	}
	return s.Start()
}

// Reset resets the stopwatch
func (s *Stopwatch) Reset() *Stopwatch {
	s.Elapsed = 0
	s.Running = false
	s.Laps = []time.Duration{}
	return s
}

// Lap records a lap time
func (s *Stopwatch) Lap() *Stopwatch {
	s.Laps = append(s.Laps, s.GetElapsed())
	return s
}

// GetElapsed returns current elapsed time
func (s *Stopwatch) GetElapsed() time.Duration {
	if s.Running {
		return s.Elapsed + time.Since(s.startTime)
	}
	return s.Elapsed
}

// Update handles messages
func (s *Stopwatch) Update(msg tea.Msg) (*Stopwatch, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "enter":
			s.Toggle()
		case "r":
			s.Reset()
		case "l":
			if s.Running {
				s.Lap()
			}
		}
	}
	return s, nil
}

// View renders the stopwatch
func (s *Stopwatch) View() string {
	elapsed := s.GetElapsed()
	hours := int(elapsed.Hours())
	minutes := int(elapsed.Minutes()) % 60
	seconds := int(elapsed.Seconds()) % 60
	ms := (elapsed.Milliseconds() % 1000) / 10

	timeStr := fmt.Sprintf("%02d:%02d:%02d.%02d", hours, minutes, seconds, ms)

	status := "⏸"
	if s.Running {
		status = "▶"
	}

	return s.style.Render(status + " " + timeStr)
}

