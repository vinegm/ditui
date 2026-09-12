package morse

import (
	"strings"
	"unicode/utf8"
)

// Elements as the player types them.
const (
	Dit = '.'
	Dah = '-'
)

// Elements as they are drawn on screen.
const (
	prettyDit = "•"
	prettyDah = "−"
)

// Longest code in the table, the one for "$".
const MaxElements = 7

// International Morse code, plus the common prosign-free punctuation.
var codes = map[rune]string{
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".", 'F': "..-.",
	'G': "--.", 'H': "....", 'I': "..", 'J': ".---", 'K': "-.-", 'L': ".-..",
	'M': "--", 'N': "-.", 'O': "---", 'P': ".--.", 'Q': "--.-", 'R': ".-.",
	'S': "...", 'T': "-", 'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-",
	'Y': "-.--", 'Z': "--..",

	'0': "-----", '1': ".----", '2': "..---", '3': "...--", '4': "....-",
	'5': ".....", '6': "-....", '7': "--...", '8': "---..", '9': "----.",

	'.': ".-.-.-", ',': "--..--", '?': "..--..", '\'': ".----.", '!': "-.-.--",
	'/': "-..-.", '(': "-.--.", ')': "-.--.-", '&': ".-...", ':': "---...",
	';': "-.-.-.", '=': "-...-", '+': ".-.-.", '-': "-....-", '_': "..--.-",
	'"': ".-..-.", '@': ".--.-.", '$': "...-..-",
}

// Reverse of codes, built once.
var chars = func() map[string]rune {
	reverse := make(map[string]rune, len(codes))

	for char, code := range codes {
		reverse[code] = char
	}

	return reverse
}()

// Code returns the dits and dahs for a character, and false if it has none.
// Letters are matched regardless of case.
func Code(char rune) (string, bool) {
	code, ok := codes[normalize(char)]
	return code, ok
}

// Lookup returns the character a code stands for, and false if none does.
func Lookup(code string) (rune, bool) {
	char, ok := chars[code]
	return char, ok
}

// Pretty draws a code with proper dot and dash glyphs: "•−••".
func Pretty(code string) string {
	return strings.NewReplacer(string(Dit), prettyDit, string(Dah), prettyDah).Replace(code)
}

// Spaced draws a code with a gap between elements, for reading at a glance: "• − • •".
func Spaced(code string) string {
	elements := make([]string, 0, len(code))

	for _, element := range Pretty(code) {
		elements = append(elements, string(element))
	}

	return strings.Join(elements, " ")
}

// ParseChar returns the character a decode answer names, and false if text
// is not exactly one character with a code.
func ParseChar(text string) (rune, bool) {
	if utf8.RuneCountInString(text) != 1 {
		return 0, false
	}

	char, _ := utf8.DecodeRuneInString(text)
	char = normalize(char)
	_, ok := codes[char]

	return char, ok
}

// IsElement reports whether text is a single dit or dah key.
func IsElement(text string) bool {
	return text == string(Dit) || text == string(Dah)
}

func normalize(char rune) rune {
	if char >= 'a' && char <= 'z' {
		return char - 'a' + 'A'
	}

	return char
}
