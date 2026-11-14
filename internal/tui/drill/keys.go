package drill

import "charm.land/bubbles/v2/key"

// keyMap is the drill's keys.
type keyMap struct {
	Finish, Answer, Delete, Clear key.Binding

	// Dit and Dah are typed into the answer, and AnyChar
	// stands for every character, all only named
	// in the key hints bar.
	Dit, Dah, AnyChar key.Binding
}

var keys = keyMap{
	Finish: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Finish"),
	),
	Answer: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter", "Answer"),
	),
	Delete: key.NewBinding(key.WithKeys("backspace")),
	Clear:  key.NewBinding(key.WithKeys("ctrl+u")),
	Dit: key.NewBinding(
		key.WithKeys("."),
		key.WithHelp(".", "Dit"),
	),
	Dah: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "Dah"),
	),
	AnyChar: key.NewBinding(
		key.WithKeys("a character"),
		key.WithHelp("a character", "Answer"),
	),
}
