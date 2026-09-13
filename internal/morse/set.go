package morse

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
)

// Name of a set given as a list of characters rather than picked by name.
const CustomName = "custom"

// Set is the pool of characters a drill asks from.
type Set struct {
	Name  string
	Chars []rune
}

// Presets are the named sets, in the order the menu cycles through them.
func Presets() []Set {
	return []Set{
		{"letters", []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
		{"digits", []rune("0123456789")},
		{"alphanumeric", []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")},
		{"punctuation", []rune(".,?'!/()&:;=+-_\"@$")},
	}
}

// Preset returns the named set, and false if no preset has that name.
func Preset(name string) (Set, bool) {
	for _, set := range Presets() {
		if set.Name == name {
			return set, true
		}
	}

	return Set{}, false
}

// Custom builds a set out of the given characters, ignoring case, spaces and
// repeats: "eish 5" drills E, I, S, H and 5.
func Custom(text string) (Set, error) {
	var chars []rune

	for _, char := range text {
		if char == ' ' {
			continue
		}

		if _, ok := Code(char); !ok {
			return Set{}, fmt.Errorf("%q has no morse code", char)
		}

		if char = normalize(char); !slices.Contains(chars, char) {
			chars = append(chars, char)
		}
	}

	if len(chars) == 0 {
		return Set{}, errors.New("no characters to drill")
	}

	return Set{Name: CustomName, Chars: chars}, nil
}

// ResolveSet picks the set a drill runs with: the custom characters when
// there are any, the named preset otherwise.
func ResolveSet(name, chars string) (Set, error) {
	if chars != "" {
		return Custom(chars)
	}

	if set, ok := Preset(name); ok {
		return set, nil
	}

	return Set{}, fmt.Errorf("unknown set %q, want one of %s", name, strings.Join(presetNames(), ", "))
}

// Label is the set as shown to the player: its name, or its characters when custom.
func (set Set) Label() string {
	if set.Name == CustomName {
		return string(set.Chars)
	}

	return set.Name
}

// Random picks a character from the set, never the previous one unless
// it is all the set has.
func (set Set) Random(previous rune) rune {
	for {
		char := set.Chars[rand.IntN(len(set.Chars))]

		if char != previous || len(set.Chars) == 1 {
			return char
		}
	}
}

// Cheatsheet lays out every character of the set next to its
// code, as many to a line as fit in width.
func (set Set) Cheatsheet(width int) []string {
	longest := 0
	for _, char := range set.Chars {
		longest = max(longest, len(codes[char]))
	}

	const gap = 2
	entryWidth := len("X ") + longest
	perLine := max(1, (width+gap)/(entryWidth+gap))

	var lines []string

	for start := 0; start < len(set.Chars); start += perLine {
		entries := make([]string, 0, perLine)

		for _, char := range set.Chars[start:min(len(set.Chars), start+perLine)] {
			code := codes[char]
			entries = append(entries, string(char)+" "+Pretty(code)+strings.Repeat(" ", longest-len(code)))
		}

		lines = append(lines, strings.TrimRight(strings.Join(entries, strings.Repeat(" ", gap)), " "))
	}

	return lines
}

func presetNames() []string {
	var names []string

	for _, set := range Presets() {
		names = append(names, set.Name)
	}

	return names
}
