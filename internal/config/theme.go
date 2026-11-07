package config

// Theme is every colour the windows render with.
type Theme struct {
	Fg       string `yaml:"fg"`
	Bg       string `yaml:"bg"`
	Selected string `yaml:"selected"`
	Border   string `yaml:"border"`
	Accent   string `yaml:"accent"` // the code or character being asked
	Muted    string `yaml:"muted"`  // hints, timers and help lines
	Correct  string `yaml:"correct"`
	Wrong    string `yaml:"wrong"`
}
