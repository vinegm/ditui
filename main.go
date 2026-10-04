package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/vinegm/ditui/internal/config"
	"github.com/vinegm/ditui/internal/morse"
	"github.com/vinegm/ditui/internal/tui"
	"github.com/vinegm/ditui/internal/tui/drill"
)

const (
	logFile  = "/tmp/ditui.log"
	logLevel = "debug"
)

var version = "dev"

type flags struct {
	modeName    string
	setName     string
	chars       string
	configPath  string
	debugFile   string
	cheatsheet  bool
	debug       bool
	showVersion bool
	noConfig    bool
}

func main() {
	opts := parseFlags()

	if opts.showVersion {
		fmt.Println("ditui version", version)
		return
	}

	log.SetOutput(io.Discard)
	if opts.debug {
		file, err := tea.LogToFile(opts.debugFile, logLevel)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	cfg := loadConfig(opts.noConfig, opts.configPath)

	program := tea.NewProgram(tui.New(cfg, drillSettings(opts, cfg)))
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() flags {
	var opts flags

	defaults := config.Default()

	flag.StringVar(&opts.modeName, "mode", defaults.Mode, "`MODE` to drill: decode (code to character) or encode (character to code)")
	flag.StringVar(&opts.setName, "set", defaults.Set, "`SET` of characters to drill: letters, digits, alphanumeric or punctuation")
	flag.StringVar(&opts.chars, "chars", defaults.Chars, "drill only these `CHARACTERS`, overriding --set")
	flag.BoolVar(&opts.cheatsheet, "cheatsheet", defaults.Cheatsheet, "show every code of the set under the prompt")
	flag.BoolVar(&opts.debug, "debug", false, "enable debug mode")
	flag.StringVar(&opts.debugFile, "debug-file", logFile, "path to debug log file (ignored if --debug is not set)")
	flag.BoolVar(&opts.showVersion, "version", false, "print the version and exit")
	flag.BoolVar(&opts.noConfig, "no-config", false, "ignore config.yaml and run with the built-in defaults")
	flag.StringVar(&opts.configPath, "config", "", "config file to load instead of the default path")
	flag.Parse()

	return opts
}

// applyConfig takes the drill settings not given on the command line from
// cfg. The config is only known after parsing, so it can't be the flags' default.
func (opts *flags) applyConfig(cfg config.Config) {
	if !isSet("mode") {
		opts.modeName = cfg.Mode
	}

	if !isSet("set") {
		opts.setName = cfg.Set
	}

	// A preset picked on the command line beats custom characters from the config.
	if !isSet("set") && !isSet("chars") {
		opts.chars = cfg.Chars
	}

	if !isSet("cheatsheet") {
		opts.cheatsheet = cfg.Cheatsheet
	}
}

// drillSettings is what a drill starts from: the
// command line, then the config for what it left out.
func drillSettings(opts flags, cfg config.Config) drill.Options {
	opts.applyConfig(validDrillConfig(cfg))

	mode, err := morse.ParseMode(opts.modeName)
	if err != nil {
		usageError(err.Error())
	}

	set, err := morse.ResolveSet(opts.setName, opts.chars)
	if err != nil {
		usageError(err.Error())
	}

	return drill.Options{
		Mode:       mode,
		Set:        set,
		Cheatsheet: opts.cheatsheet,
	}
}

// validDrillConfig resets each drill setting no drill can run
// with to its default, keeping the rest of the file.
func validDrillConfig(cfg config.Config) config.Config {
	defaults := config.Default()

	if _, err := morse.ParseMode(cfg.Mode); err != nil {
		warn("invalid mode in the config, using the default: %v", err)

		cfg.Mode = defaults.Mode
	}

	if _, err := morse.ResolveSet(cfg.Set, cfg.Chars); err != nil {
		warn("invalid set or chars in the config, using the defaults: %v", err)

		cfg.Set = defaults.Set
		cfg.Chars = defaults.Chars
	}

	return cfg
}

// loadConfig returns the defaults with -no-config, else the file at
// path or the default path. Only an explicit path must exist.
func loadConfig(noConfig bool, path string) config.Config {
	if noConfig {
		return config.Default()
	}

	cfg, err := config.Load(cmp.Or(path, config.DefaultPath()))
	if err != nil && (!errors.Is(err, fs.ErrNotExist) || path != "") {
		warn("could not load config: %v", err)
	}

	return cfg
}

// warn reports a problem found before the TUI starts, on
// stderr so it is seen without -debug, and in the debug log.
func warn(format string, args ...any) {
	text := fmt.Sprintf(format, args...)

	fmt.Fprintln(os.Stderr, "warning: "+text)
	log.Print(text)
}

// isSet reports whether the named flag was given on the command line.
func isSet(name string) bool {
	given := false

	flag.Visit(func(set *flag.Flag) {
		given = given || set.Name == name
	})

	return given
}

// usageError reports a bad command line the way the flag package does.
func usageError(message string) {
	fmt.Fprintln(flag.CommandLine.Output(), message)
	flag.Usage()
	os.Exit(2)
}
