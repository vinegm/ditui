package config

// Config is config.yaml.
// Keys are documented in Themes/example-config.yaml.
type Config struct {
	Drill Drill `yaml:"drill"`
	Theme Theme `yaml:"theme"`
}

// Drill is what a drill asks and shows.
type Drill struct {
	Mode       string `yaml:"mode"`       // "decode" (code to character) or "encode" (character to code)
	Set        string `yaml:"set"`        // letters, digits, alphanumeric or punctuation
	Chars      string `yaml:"chars"`      // drill just these characters, overriding set
	Cheatsheet bool   `yaml:"cheatsheet"` // show every code of the set under the prompt
}

// Default is the built-in configuration, used
// as the base the config file is layered over.
func Default() Config {
	return Config{
		Drill: Drill{
			Mode:       "decode",
			Set:        "letters",
			Chars:      "",
			Cheatsheet: false,
		},
		Theme: Theme{
			Fg:       "",
			Bg:       "",
			Selected: "#5a5a5a",
			Border:   "",
			Accent:   "",
			Muted:    "#6c6c6c",
			Correct:  "#5faf5f",
			Wrong:    "#d75f5f",
		},
	}
}
