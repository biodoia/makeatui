// Package advanced - Calendar and Clock widgets
package advanced

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Calendar provides a calendar widget
type Calendar struct {
	ID           string
	Year         int
	Month        time.Month
	SelectedDay  int
	Today        time.Time
	MarkedDays   map[int]string // day -> marker color
	FirstDayMon  bool           // Week starts on Monday
	style        CalendarStyle
}

// CalendarStyle holds styling
type CalendarStyle struct {
	Container   lipgloss.Style
	Header      lipgloss.Style
	DayName     lipgloss.Style
	Day         lipgloss.Style
	DayToday    lipgloss.Style
	DaySelected lipgloss.Style
	DayMarked   lipgloss.Style
	DayOther    lipgloss.Style
	NavButton   lipgloss.Style
}

// DefaultCalendarStyle returns default styling
func DefaultCalendarStyle() CalendarStyle {
	return CalendarStyle{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B2CBF")).
			Padding(1),
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true).
			Align(lipgloss.Center),
		DayName: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9D4EDD")).
			Bold(true).
			Width(4).
			Align(lipgloss.Center),
		Day: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Width(4).
			Align(lipgloss.Center),
		DayToday: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E040FB")).
			Bold(true).
			Width(4).
			Align(lipgloss.Center),
		DaySelected: lipgloss.NewStyle().
			Background(lipgloss.Color("#9D4EDD")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Width(4).
			Align(lipgloss.Center),
		DayMarked: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6BCB77")).
			Width(4).
			Align(lipgloss.Center),
		DayOther: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Width(4).
			Align(lipgloss.Center),
		NavButton: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7B2CBF")),
	}
}

// NewCalendar creates a calendar widget
func NewCalendar(id string) *Calendar {
	now := time.Now()
	return &Calendar{
		ID:          id,
		Year:        now.Year(),
		Month:       now.Month(),
		SelectedDay: now.Day(),
		Today:       now,
		MarkedDays:  make(map[int]string),
		FirstDayMon: true,
		style:       DefaultCalendarStyle(),
	}
}

// SetDate sets the calendar date
func (c *Calendar) SetDate(year int, month time.Month, day int) *Calendar {
	c.Year = year
	c.Month = month
	c.SelectedDay = day
	return c
}

// MarkDay marks a day with a color
func (c *Calendar) MarkDay(day int, color string) *Calendar {
	c.MarkedDays[day] = color
	return c
}

// NextMonth moves to next month
func (c *Calendar) NextMonth() *Calendar {
	c.Month++
	if c.Month > 12 {
		c.Month = 1
		c.Year++
	}
	return c
}

// PrevMonth moves to previous month
func (c *Calendar) PrevMonth() *Calendar {
	c.Month--
	if c.Month < 1 {
		c.Month = 12
		c.Year--
	}
	return c
}

// daysInMonth returns days in current month
func (c *Calendar) daysInMonth() int {
	return time.Date(c.Year, c.Month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// firstWeekday returns weekday of first day (0=Sunday, 1=Monday, etc.)
func (c *Calendar) firstWeekday() int {
	d := time.Date(c.Year, c.Month, 1, 0, 0, 0, 0, time.UTC)
	wd := int(d.Weekday())
	if c.FirstDayMon {
		wd = (wd + 6) % 7
	}
	return wd
}

// Update handles messages
func (c *Calendar) Update(msg tea.Msg) (*Calendar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			c.SelectedDay--
			if c.SelectedDay < 1 {
				c.PrevMonth()
				c.SelectedDay = c.daysInMonth()
			}
		case "right", "l":
			c.SelectedDay++
			if c.SelectedDay > c.daysInMonth() {
				c.NextMonth()
				c.SelectedDay = 1
			}
		case "up", "k":
			c.SelectedDay -= 7
			if c.SelectedDay < 1 {
				c.PrevMonth()
				c.SelectedDay += c.daysInMonth()
			}
		case "down", "j":
			c.SelectedDay += 7
			if c.SelectedDay > c.daysInMonth() {
				c.SelectedDay -= c.daysInMonth()
				c.NextMonth()
			}
		case "[":
			c.PrevMonth()
		case "]":
			c.NextMonth()
		case "t":
			c.Year = c.Today.Year()
			c.Month = c.Today.Month()
			c.SelectedDay = c.Today.Day()
		}
	}
	return c, nil
}

// View renders the calendar
func (c *Calendar) View() string {
	var lines []string

	// Header
	header := c.style.NavButton.Render("◀ ") +
		c.style.Header.Render(fmt.Sprintf("%s %d", c.Month.String(), c.Year)) +
		c.style.NavButton.Render(" ▶")
	lines = append(lines, header)
	lines = append(lines, "")

	// Day names
	dayNames := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	if !c.FirstDayMon {
		dayNames = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	}
	var dayNameRow []string
	for _, d := range dayNames {
		dayNameRow = append(dayNameRow, c.style.DayName.Render(d))
	}
	lines = append(lines, strings.Join(dayNameRow, ""))

	// Days grid
	firstWd := c.firstWeekday()
	daysInMo := c.daysInMonth()
	day := 1

	for week := 0; week < 6; week++ {
		var weekRow []string
		for wd := 0; wd < 7; wd++ {
			if (week == 0 && wd < firstWd) || day > daysInMo {
				weekRow = append(weekRow, c.style.DayOther.Render("  "))
			} else {
				style := c.style.Day
				if day == c.SelectedDay {
					style = c.style.DaySelected
				} else if c.Year == c.Today.Year() && c.Month == c.Today.Month() && day == c.Today.Day() {
					style = c.style.DayToday
				} else if _, marked := c.MarkedDays[day]; marked {
					style = c.style.DayMarked
				}

				dayStr := fmt.Sprintf("%2d", day)
				weekRow = append(weekRow, style.Render(dayStr))
				day++
			}
		}
		lines = append(lines, strings.Join(weekRow, ""))
		if day > daysInMo {
			break
		}
	}

	return c.style.Container.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// GetSelectedDate returns the selected date
func (c *Calendar) GetSelectedDate() time.Time {
	return time.Date(c.Year, c.Month, c.SelectedDay, 0, 0, 0, 0, time.UTC)
}

