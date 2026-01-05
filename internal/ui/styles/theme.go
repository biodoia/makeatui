// Package styles provides glamorous theming for MakeaTUI
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines the color palette and styles for the TUI
type Theme struct {
	Name string

	// Primary colors - Ultraviolet palette
	Primary      lipgloss.Color
	Secondary    lipgloss.Color
	Accent       lipgloss.Color
	Background   lipgloss.Color
	Surface      lipgloss.Color
	SurfaceLight lipgloss.Color

	// Text colors
	TextPrimary   lipgloss.Color
	TextSecondary lipgloss.Color
	TextMuted     lipgloss.Color

	// Semantic colors
	Success lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color
	Info    lipgloss.Color

	// Border colors
	Border       lipgloss.Color
	BorderActive lipgloss.Color
}

// Ultraviolet is the default glamorous theme inspired by Charm
var Ultraviolet = Theme{
	Name:         "Ultraviolet",
	Primary:      lipgloss.Color("#9D4EDD"), // Vivid purple
	Secondary:    lipgloss.Color("#7B2CBF"), // Deep purple
	Accent:       lipgloss.Color("#E040FB"), // Pink/magenta
	Background:   lipgloss.Color("#0D0221"), // Deep space
	Surface:      lipgloss.Color("#1A1333"), // Slightly lighter
	SurfaceLight: lipgloss.Color("#2D1B4E"), // Elevated surface

	TextPrimary:   lipgloss.Color("#FFFFFF"),
	TextSecondary: lipgloss.Color("#B8B8D1"),
	TextMuted:     lipgloss.Color("#6B6B8D"),

	Success: lipgloss.Color("#00E676"),
	Warning: lipgloss.Color("#FFD600"),
	Error:   lipgloss.Color("#FF5252"),
	Info:    lipgloss.Color("#40C4FF"),

	Border:       lipgloss.Color("#3D2C5E"),
	BorderActive: lipgloss.Color("#9D4EDD"),
}

// Neon is an alternative high-contrast neon theme
var Neon = Theme{
	Name:         "Neon",
	Primary:      lipgloss.Color("#00FFFF"), // Cyan
	Secondary:    lipgloss.Color("#FF00FF"), // Magenta
	Accent:       lipgloss.Color("#FFFF00"), // Yellow
	Background:   lipgloss.Color("#000000"),
	Surface:      lipgloss.Color("#0A0A0A"),
	SurfaceLight: lipgloss.Color("#1A1A1A"),

	TextPrimary:   lipgloss.Color("#FFFFFF"),
	TextSecondary: lipgloss.Color("#AAAAAA"),
	TextMuted:     lipgloss.Color("#666666"),

	Success: lipgloss.Color("#00FF00"),
	Warning: lipgloss.Color("#FFAA00"),
	Error:   lipgloss.Color("#FF0000"),
	Info:    lipgloss.Color("#00AAFF"),

	Border:       lipgloss.Color("#333333"),
	BorderActive: lipgloss.Color("#00FFFF"),
}

// CurrentTheme holds the active theme
var CurrentTheme = Ultraviolet

// Styles contains all the pre-built styles for the application
type Styles struct {
	// Layout styles
	App       lipgloss.Style
	Sidebar   lipgloss.Style
	Canvas    lipgloss.Style
	Toolbar   lipgloss.Style
	StatusBar lipgloss.Style

	// Component styles
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	Label         lipgloss.Style
	Button        lipgloss.Style
	ButtonActive  lipgloss.Style
	Input         lipgloss.Style
	InputFocused  lipgloss.Style
	ListItem      lipgloss.Style
	ListItemActive lipgloss.Style
	Box           lipgloss.Style
	BoxSelected   lipgloss.Style

	// Special styles
	Logo    lipgloss.Style
	Help    lipgloss.Style
	Error   lipgloss.Style
	Success lipgloss.Style
}

// NewStyles creates a new Styles instance based on the given theme
func NewStyles(t Theme) Styles {
	return Styles{
		App: lipgloss.NewStyle().
			Background(t.Background),

		Sidebar: lipgloss.NewStyle().
			Background(t.Surface).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Border).
			Padding(1, 2),

		Canvas: lipgloss.NewStyle().
			Background(t.Background).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary).
			Padding(1, 2),

		Toolbar: lipgloss.NewStyle().
			Background(t.SurfaceLight).
			Foreground(t.TextSecondary).
			Padding(0, 2),

		StatusBar: lipgloss.NewStyle().
			Background(t.Primary).
			Foreground(t.TextPrimary).
			Padding(0, 2),

		Title: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(t.TextSecondary).
			Italic(true),

		Label: lipgloss.NewStyle().
			Foreground(t.TextSecondary),

		Button: lipgloss.NewStyle().
			Background(t.Surface).
			Foreground(t.TextPrimary).
			Padding(0, 3).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Border),

		ButtonActive: lipgloss.NewStyle().
			Background(t.Primary).
			Foreground(t.TextPrimary).
			Padding(0, 3).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent).
			Bold(true),

		Input: lipgloss.NewStyle().
			Background(t.Surface).
			Foreground(t.TextPrimary).
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Border),

		InputFocused: lipgloss.NewStyle().
			Background(t.Surface).
			Foreground(t.TextPrimary).
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary),

		ListItem: lipgloss.NewStyle().
			Foreground(t.TextSecondary).
			PaddingLeft(2),

		ListItemActive: lipgloss.NewStyle().
			Foreground(t.TextPrimary).
			Background(t.SurfaceLight).
			Bold(true).
			PaddingLeft(1),

		Box: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Border).
			Padding(1, 2),

		BoxSelected: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(t.Primary).
			Padding(1, 2),

		Logo: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),

		Help: lipgloss.NewStyle().
			Foreground(t.TextMuted),

		Error: lipgloss.NewStyle().
			Foreground(t.Error).
			Bold(true),

		Success: lipgloss.NewStyle().
			Foreground(t.Success).
			Bold(true),
	}
}

