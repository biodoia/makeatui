// Package animation provides spring-based animations using Harmonica
package animation

import (
	"time"

	"github.com/charmbracelet/harmonica"
)

// Spring wraps harmonica's spring animation for smooth transitions
type Spring struct {
	spring   harmonica.Spring
	position float64
	velocity float64
	target   float64
}

// SpringConfig defines spring animation parameters
type SpringConfig struct {
	Frequency    float64 // Oscillation frequency (Hz)
	DampingRatio float64 // Damping ratio (0-1, 1 = critical damping)
}

// DefaultSpring returns a smooth, responsive spring config
var DefaultSpring = SpringConfig{
	Frequency:    7.0,
	DampingRatio: 0.8,
}

// SmoothSpring returns a slower, more fluid spring
var SmoothSpring = SpringConfig{
	Frequency:    4.0,
	DampingRatio: 0.9,
}

// BouncySpring returns a springy, bouncy animation
var BouncySpring = SpringConfig{
	Frequency:    10.0,
	DampingRatio: 0.5,
}

// NewSpring creates a new spring animation
func NewSpring(config SpringConfig) *Spring {
	return &Spring{
		spring: harmonica.NewSpring(harmonica.FPS(60), config.Frequency, config.DampingRatio),
	}
}

// SetTarget sets the target value for the spring
func (s *Spring) SetTarget(target float64) {
	s.target = target
}

// Update advances the spring animation by one frame
func (s *Spring) Update() float64 {
	s.position, s.velocity = s.spring.Update(s.position, s.velocity, s.target)
	return s.position
}

// Position returns the current position
func (s *Spring) Position() float64 {
	return s.position
}

// IsSettled returns true if the spring has settled at the target
func (s *Spring) IsSettled() bool {
	const threshold = 0.01
	return abs(s.position-s.target) < threshold && abs(s.velocity) < threshold
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Animator manages multiple spring animations
type Animator struct {
	springs   map[string]*Spring
	ticker    *time.Ticker
	running   bool
	frameRate time.Duration
}

// NewAnimator creates a new animation manager
func NewAnimator(fps int) *Animator {
	return &Animator{
		springs:   make(map[string]*Spring),
		frameRate: time.Second / time.Duration(fps),
	}
}

// AddSpring adds a named spring animation
func (a *Animator) AddSpring(name string, config SpringConfig) *Spring {
	spring := NewSpring(config)
	a.springs[name] = spring
	return spring
}

// GetSpring returns a spring by name
func (a *Animator) GetSpring(name string) *Spring {
	return a.springs[name]
}

// UpdateAll updates all springs and returns true if any are still animating
func (a *Animator) UpdateAll() bool {
	anyActive := false
	for _, spring := range a.springs {
		spring.Update()
		if !spring.IsSettled() {
			anyActive = true
		}
	}
	return anyActive
}

// AnimateValue is a convenience function for animating a single value
func AnimateValue(from, to float64, duration time.Duration, callback func(float64)) {
	spring := NewSpring(DefaultSpring)
	spring.position = from
	spring.SetTarget(to)

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	timeout := time.After(duration)

	for {
		select {
		case <-ticker.C:
			pos := spring.Update()
			callback(pos)
			if spring.IsSettled() {
				callback(to) // Ensure we end exactly at target
				return
			}
		case <-timeout:
			callback(to)
			return
		}
	}
}

// EaseInOut provides a smooth easing curve
func EaseInOut(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - (-2*t+2)*(-2*t+2)/2
}

// Lerp linearly interpolates between two values
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

