package config

type Config struct {
	// Drill settings
	Mode       string `yaml:"mode"`       // "decode" (code to character) or "encode" (character to code)
	Set        string `yaml:"set"`        // letters, digits, alphanumeric or punctuation
	Chars      string `yaml:"chars"`      // drill just these characters, overriding set
	Cheatsheet bool   `yaml:"cheatsheet"` // show every code of the set under the prompt

	// Style settings
	FgColor       string `yaml:"fg_color"`
	BgColor       string `yaml:"bg_color"`
	SelectedColor string `yaml:"selected_color"`
	BorderColor   string `yaml:"border_color"`
	AccentColor   string `yaml:"accent_color"` // the code or character being asked
	MutedColor    string `yaml:"muted_color"`  // hints, timers and help lines
	CorrectColor  string `yaml:"correct_color"`
	WrongColor    string `yaml:"wrong_color"`
}

// The built-in configuration, used as the
// base the config file is layered over.
func Default() Config {
	return Config{
		Mode:       "decode",
		Set:        "letters",
		Chars:      "",
		Cheatsheet: false,

		FgColor:       "",
		BgColor:       "",
		SelectedColor: "#5a5a5a",
		BorderColor:   "",
		AccentColor:   "",
		MutedColor:    "#6c6c6c",
		CorrectColor:  "#5faf5f",
		WrongColor:    "#d75f5f",
	}
}
