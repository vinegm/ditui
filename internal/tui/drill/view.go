package drill

import (
	"fmt"
	"strings"
	"time"

	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/ditui/internal/morse"
	"github.com/vinegm/ditui/internal/report"
)

// How many of the latest answers the history strip shows.
const historyLength = 25

// MaxCheatsheetLines is the most lines the cheatsheet takes, so a large custom set can't outgrow the terminal.
const MaxCheatsheetLines = 6

const (
	decodeHelp = "type the character · esc finish"
	encodeHelp = ". dit · - dah · enter answer · esc finish"
)

// SetSize does nothing, the drill is drawn at its own size.
func (window *Window) SetSize(_, _ int) {}

// View renders the drill. Every line is always drawn, blank if unused,
// so the window keeps its height from one question to the next.
func (window *Window) View() string {
	th := window.theme

	lines := []string{
		th.TitleStyle.Render("ditui") + th.MutedStyle.Render("  "+window.options.Mode.String()+" · "+window.options.Set.Label()),
		th.TextStyle.Render(window.score()),
		window.history(),
		"",
		th.AccentStyle.Render(window.question()),
		"",
		window.inputLine(),
		th.MutedStyle.Render(fmt.Sprintf("%.1fs", time.Since(window.askedAt).Seconds())),
		"",
	}

	lines = append(lines, window.feedbackLines()...)
	lines = append(lines, "")

	if window.options.Cheatsheet {
		// Left-aligned so the columns line up and center as a block.
		lines = append(lines, th.MutedStyle.Align(lg.Left).Render(strings.Join(window.cheatsheet(), "\n")), "")
	}

	lines = append(lines, th.MutedStyle.Render(window.help()))

	return th.WindowStyle.Render(strings.Join(lines, "\n"))
}

// question is what the player translates: the code when
// decoding, the character when encoding.
func (window *Window) question() string {
	if window.options.Mode == morse.Encode {
		return string(window.char)
	}

	code, _ := morse.Code(window.char)

	return morse.Spaced(code)
}

func (window *Window) help() string {
	if window.options.Mode == morse.Encode {
		return encodeHelp
	}

	return decodeHelp
}

func (window *Window) score() string {
	if len(window.attempts) == 0 {
		return "no answers yet"
	}

	correct := report.CountCorrect(window.attempts)

	return fmt.Sprintf("%d/%d correct (%.0f%%) · avg %.1fs",
		correct, len(window.attempts), 100*float64(correct)/float64(len(window.attempts)), report.Mean(window.attempts))
}

// history is a strip of dots, one per recent answer, colored by whether it was right.
func (window *Window) history() string {
	recent := window.attempts[max(0, len(window.attempts)-historyLength):]

	var sb strings.Builder

	for _, attempt := range recent {
		if attempt.Correct() {
			sb.WriteString(window.theme.CorrectStyle.Render("●"))
		} else {
			sb.WriteString(window.theme.WrongStyle.Render("●"))
		}
	}

	return sb.String()
}

func (window *Window) inputLine() string {
	prompt := window.theme.SelectedMenuTextStyle.Render("› ")

	if window.options.Mode == morse.Decode {
		return prompt + window.theme.MutedStyle.Render("which character?")
	}

	if window.input == "" {
		return prompt + window.theme.MutedStyle.Render("key the code")
	}

	return prompt + window.theme.TextStyle.Render(morse.Spaced(window.input)) + window.theme.SelectedMenuTextStyle.Render("▏")
}

// feedbackLines are the verdict on the previous submission and its detail.
func (window *Window) feedbackLines() []string {
	last := window.last
	if last == nil {
		return []string{"", ""}
	}

	th := window.theme
	verdict := lg.NewStyle().MaxWidth(th.InnerWidth()).Render(last.String())

	style := th.WrongStyle
	if last.outcome == answeredRight {
		style = th.CorrectStyle
	}

	return []string{style.Render(verdict), th.MutedStyle.Render(last.Detail())}
}

// cheatsheet is the set's codes, cut short with a count of what didn't fit.
func (window *Window) cheatsheet() []string {
	lines := window.options.Set.Cheatsheet(window.theme.InnerWidth())

	if len(lines) > MaxCheatsheetLines {
		lines = append(lines[:MaxCheatsheetLines-1], "…")
	}

	return lines
}
