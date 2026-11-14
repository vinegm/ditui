package menu

import "charm.land/bubbles/v2/key"

// keyMap is the menu's keys.
type keyMap struct {
	Up, Down, Left, Right, Select, Quit key.Binding

	// Move and Change name the pairs of keys
	// above as one in the key hints bar.
	Move, Change key.Binding
}

var keys = keyMap{
	Up:    key.NewBinding(key.WithKeys("up", "k")),
	Down:  key.NewBinding(key.WithKeys("down", "j")),
	Left:  key.NewBinding(key.WithKeys("left", "h")),
	Right: key.NewBinding(key.WithKeys("right", "l")),
	Select: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter", "Select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Quit"),
	),
	Move: key.NewBinding(
		key.WithKeys("down", "j", "up", "k"),
		key.WithHelp("jk", "Move"),
	),
	Change: key.NewBinding(
		key.WithKeys("left", "h", "right", "l"),
		key.WithHelp("hl", "Change"),
	),
}
