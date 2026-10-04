package drill

import (
	"fmt"
	"time"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"

	"github.com/vinegm/ditui/internal/morse"
	"github.com/vinegm/ditui/internal/tui/theme"
)

// How often the question timer on screen refreshes.
const tickInterval = 100 * time.Millisecond

// Options control what the drill asks and shows.
type Options struct {
	Mode       morse.Mode
	Set        morse.Set
	Cheatsheet bool
}

type (
	// FinishedMsg asks for control to end the session, handing over what was answered.
	FinishedMsg struct{ Attempts []morse.Attempt }

	// Message carrying a timer tick, tagged with the session that armed it.
	tickMsg struct{ tag int }
)

type outcome int

const (
	answeredRight outcome = iota
	answeredWrong
	notACharacter
)

// feedback is what the previous submission turned out to be.
type feedback struct {
	outcome outcome
	attempt morse.Attempt
	input   string
}

// Window asks one character after another and times each answer.
type Window struct {
	options  Options
	attempts []morse.Attempt

	char    rune
	askedAt time.Time
	input   string // dits and dahs keyed so far, when encoding
	last    *feedback

	// Bumped every session so ticks armed by an old one are dropped.
	tickTag int

	theme *theme.Theme
}

// New initializes a new Window.
//
//	NOTE: Nothing is asked until Start() is called.
func New(th *theme.Theme) *Window {
	return &Window{theme: th}
}

// Start begins a fresh session with the given options.
func (window *Window) Start(options Options) tea.Cmd {
	window.options = options
	window.attempts = nil
	window.last = nil
	window.char = 0
	window.tickTag++

	window.nextQuestion()

	return tick(window.tickTag)
}

// Update handles the keys and keeps the timer ticking.
func (window *Window) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return window.handleKey(msg)

	case tickMsg:
		if msg.tag != window.tickTag {
			return nil
		}

		return tick(window.tickTag)
	}

	return nil
}

func (window *Window) handleKey(msg tea.KeyPressMsg) tea.Cmd {
	if msg.String() == "esc" {
		return window.finish()
	}

	if window.options.Mode == morse.Encode {
		window.handleEncodeKey(msg)
		return nil
	}

	// Every answer is a single character, so it goes in without waiting for enter.
	if utf8.RuneCountInString(msg.Text) == 1 && msg.Text != " " {
		window.submitChar(msg.Text)
	}

	return nil
}

func (window *Window) handleEncodeKey(msg tea.KeyPressMsg) {
	switch msg.String() {
	case "enter", "space":
		window.submitCode()

	case "backspace":
		if window.input != "" {
			window.input = window.input[:len(window.input)-1]
		}

	case "ctrl+u":
		window.input = ""

	default:
		if morse.IsElement(msg.Text) && len(window.input) < morse.MaxElements {
			window.input += msg.Text
		}
	}
}

func (window *Window) submitChar(text string) {
	guess, ok := morse.ParseChar(text)
	if !ok {
		// Same code again: the clock keeps running, a slip costs time like a slow answer.
		window.last = &feedback{outcome: notACharacter, input: text}
		return
	}

	window.record(string(guess))
}

func (window *Window) submitCode() {
	if window.input == "" {
		return
	}

	window.record(window.input)
}

func (window *Window) record(guess string) {
	attempt := morse.Attempt{
		Char:    window.char,
		Mode:    window.options.Mode,
		Guess:   guess,
		Elapsed: time.Since(window.askedAt),
	}
	window.attempts = append(window.attempts, attempt)

	window.last = &feedback{outcome: answeredWrong, attempt: attempt}
	if attempt.Correct() {
		window.last.outcome = answeredRight
	}

	window.nextQuestion()
}

func (window *Window) nextQuestion() {
	window.char = window.options.Set.Random(window.char)
	window.askedAt = time.Now()
	window.input = ""
}

func (window *Window) finish() tea.Cmd {
	window.tickTag++
	attempts := window.attempts

	return func() tea.Msg {
		return FinishedMsg{Attempts: attempts}
	}
}

func tick(tag int) tea.Cmd {
	return tea.Tick(tickInterval, func(time.Time) tea.Msg {
		return tickMsg{tag: tag}
	})
}

// String form of a feedback line's verdict, used by the view.
func (last feedback) String() string {
	attempt := last.attempt
	verdict := fmt.Sprintf("%c is %s", attempt.Char, morse.Spaced(attempt.Code()))

	switch last.outcome {
	case answeredRight:
		return "✓ " + verdict

	case answeredWrong:
		return "✗ " + verdict

	default:
		return fmt.Sprintf("? %q has no morse code", last.input)
	}
}

// Detail is the line under the verdict: what was answered and how fast.
func (last feedback) Detail() string {
	attempt := last.attempt

	switch last.outcome {
	case answeredRight:
		return fmt.Sprintf("%.1fs", attempt.Seconds())

	case answeredWrong:
		return fmt.Sprintf("you said %s · %.1fs", said(attempt), attempt.Seconds())

	default:
		return "try again, the clock is still running"
	}
}

// said shows a wrong guess with what it translates
// to, so the mix-up behind it is plain to see.
func said(attempt morse.Attempt) string {
	if attempt.Mode == morse.Encode {
		if char, ok := morse.Lookup(attempt.Guess); ok {
			return fmt.Sprintf("%s (%c)", morse.Spaced(attempt.Guess), char)
		}

		return morse.Spaced(attempt.Guess) + " (nothing)"
	}

	code, _ := morse.Code([]rune(attempt.Guess)[0])

	return fmt.Sprintf("%s (%s)", attempt.Guess, morse.Spaced(code))
}
