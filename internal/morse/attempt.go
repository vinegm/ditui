package morse

import (
	"fmt"
	"time"
)

// Mode is which way a drill translates.
type Mode uint8

const (
	// Decode shows a code and asks for its character.
	Decode Mode = iota

	// Encode shows a character and asks for its code.
	Encode
)

// Modes in the order the menu cycles through them.
var Modes = []Mode{Decode, Encode}

func (mode Mode) String() string {
	if mode == Encode {
		return "encode"
	}

	return "decode"
}

// ParseMode returns the mode called name.
func ParseMode(name string) (Mode, error) {
	for _, mode := range Modes {
		if mode.String() == name {
			return mode, nil
		}
	}

	return Decode, fmt.Errorf("unknown mode %q, want decode or encode", name)
}

// Attempt is one answered question.
type Attempt struct {
	Char    rune
	Mode    Mode
	Guess   string // the character typed when decoding, the code when encoding
	Elapsed time.Duration
}

// Answer is what the guess had to be.
func (attempt Attempt) Answer() string {
	if attempt.Mode == Encode {
		return codes[attempt.Char]
	}

	return string(attempt.Char)
}

// Correct reports whether the guess was the answer.
func (attempt Attempt) Correct() bool {
	return attempt.Guess == attempt.Answer()
}

// Code is the code of the character asked.
func (attempt Attempt) Code() string {
	return codes[attempt.Char]
}

// Seconds is how long the answer took.
func (attempt Attempt) Seconds() float64 {
	return attempt.Elapsed.Seconds()
}
