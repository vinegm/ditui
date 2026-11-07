package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/ditui/internal/tui/drill"
	"github.com/vinegm/ditui/internal/tui/menu"
	"github.com/vinegm/ditui/internal/tui/results"

	"github.com/vinegm/ditui/internal/config"
	"github.com/vinegm/ditui/internal/morse"
	"github.com/vinegm/ditui/internal/report"
	"github.com/vinegm/ditui/internal/tui/theme"
)

// Smallest terminal every window and the key hints bar fit in.
const (
	minWidth  = theme.WindowWidth + 2
	minHeight = 14 + drill.MaxCheatsheetLines + hintsHeight // the drill with its cheatsheet is the tallest window
)

// Rows the key hints bar takes at the bottom of the screen.
const hintsHeight = 1

// Window is a screen of the app. SetSize is called after
// every Update, before View draws at that size.
type Window interface {
	Update(msg tea.Msg) tea.Cmd
	SetSize(width, height int)
	View() string
	Hints() []theme.Hint
}

// Model represents the state of the application.
type Model struct {
	width  int
	height int

	currentWindow Window
	menu          *menu.Window
	drill         *drill.Window
	results       *results.Window

	options drill.Options // of the drill last started, reused by "again"

	theme *theme.Theme
}

// New creates a new model with initial state.
func New(cfg config.Config, options drill.Options) Model {
	th := theme.NewTheme(cfg)
	menuWindow := menu.New(th, options)

	return Model{
		currentWindow: menuWindow,
		menu:          menuWindow,
		drill:         drill.New(th),
		results:       results.New(th),
		options:       options,
		theme:         th,
	}
}

// Init is the command run on start, none.
func (model Model) Init() tea.Cmd {
	return nil
}

// Update handles a message, then sizes the current window for View.
func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := model.update(msg)
	model.currentWindow.SetSize(model.width, model.height-hintsHeight)

	return model, cmd
}

func (model Model) update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		return model, nil

	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return model, tea.Quit
		}

	case menu.StartDrillMsg:
		model.options = msg.Options
		model.currentWindow = model.drill
		return model, model.drill.Start(msg.Options)

	case drill.FinishedMsg:
		if !model.finish(msg.Attempts) {
			model.currentWindow = model.menu
			return model, nil
		}

		model.currentWindow = model.results
		return model, nil

	case results.AgainMsg:
		model.currentWindow = model.drill
		return model, model.drill.Start(model.options)

	case results.ExitMsg:
		model.currentWindow = model.menu
		return model, nil
	}

	return model, model.currentWindow.Update(msg)
}

// finish loads the results window with a session. False when nothing was
// answered, so there is nothing to show.
func (model *Model) finish(attempts []morse.Attempt) bool {
	if len(attempts) == 0 {
		return false
	}

	model.results.Show(attempts, report.Breakdowns(attempts))

	return true
}

// View draws the current window centred on the screen, with the
// key hints bar under it, or a notice if the terminal is too small.
func (model Model) View() tea.View {
	view := tea.NewView(model.viewContent())
	view.AltScreen = true

	return view
}

func (model Model) viewContent() string {
	if model.width < minWidth || model.height < minHeight {
		return lg.Place(
			model.width,
			model.height,
			lg.Center,
			lg.Center,
			fmt.Sprintf("Terminal too small — need at least %dx%d", minWidth, minHeight),
		)
	}

	window := lg.Place(
		model.width,
		model.height-hintsHeight,
		lg.Center,
		lg.Center,
		model.currentWindow.View(),
	)

	hints := model.theme.RenderHints(model.currentWindow.Hints(), model.width)

	return lg.JoinVertical(lg.Left, window, hints)
}
