package menu

import (
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/ditui/internal/morse"
	"github.com/vinegm/ditui/internal/tui/drill"
	"github.com/vinegm/ditui/internal/tui/theme"
)

type menuOption int

const (
	startDrill menuOption = iota
	mode
	characters
	cheatsheet
	quit
)

const help = "↑/↓ move · ←/→ change · enter select · q quit"

// A single selectable row of the menu.
type menuItem struct {
	label  string
	action menuOption
}

// Asks for control to start a drill with these options.
type StartDrillMsg struct{ Options drill.Options }

// Window is the main menu, where the user picks
// the mode, characters and hints before drilling.
type Window struct {
	Choice  int
	Options drill.Options

	// Sets the characters row cycles through: the presets, then the
	// custom set the program was started with, if any.
	sets []morse.Set

	theme *theme.Theme
}

// New creates a new Window with the provided theme and options.
func New(th *theme.Theme, options drill.Options) *Window {
	sets := morse.Presets()
	if options.Set.Name == morse.CustomName {
		sets = append(sets, options.Set)
	}

	return &Window{
		Options: options,
		sets:    sets,
		theme:   th,
	}
}

func (window *Window) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return window.handleKey(msg.String())
	}

	return nil
}

// View renders the menu window as a string.
func (window *Window) View(_, _ int) string {
	var sb strings.Builder

	// Left-aligned so the art's lines are padded to one width and center as a block.
	sb.WriteString(window.theme.TitleStyle.Align(lg.Left).Render(getMenuTitle()))
	sb.WriteString("\n")
	sb.WriteString(window.theme.MutedStyle.Render("Morse code trainer"))
	sb.WriteString("\n")

	for i, item := range window.items() {
		sb.WriteString("\n")
		if window.Choice == i {
			sb.WriteString(window.theme.SelectedMenuTextStyle.Render(item.label))
			continue
		}

		sb.WriteString(window.theme.TextStyle.Render(item.label))
	}

	sb.WriteString("\n\n")
	sb.WriteString(window.theme.MutedStyle.Render(help))

	return window.theme.WindowStyle.Render(sb.String())
}

func getMenuTitle() string {
	art := `
   ▐▌ ▄   ▐▌
▗▞▀▜▌ ▄ ▗▄▟▙▄▖
▝▚▄▟▌ █   ▐▙▄▖
 − • •  • •  −`

	return art[1:] // Remove leading newline
}

func (window *Window) items() []menuItem {
	options := window.Options

	return []menuItem{
		{"Start Drill", startDrill},
		{"Mode: " + capitalize(options.Mode.String()), mode},
		{"Characters: " + capitalize(options.Set.Label()), characters},
		{"Cheatsheet: " + onOff(options.Cheatsheet), cheatsheet},
		{"Quit", quit},
	}
}

func capitalize(text string) string {
	if text == "" {
		return text
	}

	return strings.ToUpper(text[:1]) + text[1:]
}

func onOff(on bool) string {
	if on {
		return "On"
	}

	return "Off"
}

// current is the highlighted row's option.
//
//	NOTE: Relies on items() listing every option, in declaration order.
func (window *Window) current() menuOption {
	return menuOption(window.Choice)
}

func (window *Window) handleKey(key string) tea.Cmd {
	items := window.items()

	switch key {
	case "up", "k":
		if window.Choice > 0 {
			window.Choice--
		}

	case "down", "j":
		if window.Choice < len(items)-1 {
			window.Choice++
		}

	case "left", "h":
		window.adjust(-1)

	case "right", "l":
		window.adjust(1)

	case "enter", "space":
		switch window.current() {
		case startDrill:
			options := window.Options
			return func() tea.Msg {
				return StartDrillMsg{Options: options}
			}

		case mode, characters, cheatsheet:
			window.adjust(1)

		case quit:
			return tea.Quit
		}

	case "q":
		return tea.Quit
	}

	return nil
}

// adjust steps the setting the highlighted row edits.
func (window *Window) adjust(step int) {
	options := &window.Options

	switch window.current() {
	case mode:
		options.Mode = cycle(morse.Modes, options.Mode, step)

	case characters:
		index := slices.IndexFunc(window.sets, func(set morse.Set) bool {
			return set.Name == options.Set.Name
		})
		options.Set = window.sets[wrap(index+step, len(window.sets))]

	case cheatsheet:
		options.Cheatsheet = !options.Cheatsheet
	}
}

// cycle returns the value step places after current in values, wrapping around.
func cycle[T comparable](values []T, current T, step int) T {
	return values[wrap(slices.Index(values, current)+step, len(values))]
}

func wrap(index, length int) int {
	return ((index % length) + length) % length
}
