package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"

	"dit/internal/config"
	"dit/internal/morse"
	"dit/internal/tui"
	"dit/internal/tui/drill"
)

const (
	logFile  = "/tmp/dit.log"
	logLevel = "debug"
)

var version = "dev"

type flags struct {
	modeName    string
	setName     string
	chars       string
	debugFile   string
	cheatsheet  bool
	debug       bool
	showVersion bool
}

func main() {
	// Loaded first so the flags default to it and only override what they set.
	cfg, cfgErr := config.Load(config.DefaultPath())

	opts := parseFlags(cfg)

	if opts.showVersion {
		fmt.Println("dit version", version)
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

	if cfgErr != nil && !os.IsNotExist(cfgErr) {
		log.Printf("Could not load config: %v", cfgErr)
	}

	mode, err := morse.ParseMode(opts.modeName)
	if err != nil {
		usageError(err.Error())
	}

	set, err := morse.ResolveSet(opts.setName, opts.chars)
	if err != nil {
		usageError(err.Error())
	}

	options := drill.Options{
		Mode:       mode,
		Set:        set,
		Cheatsheet: opts.cheatsheet,
	}

	program := tea.NewProgram(tui.New(cfg, options))
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseFlags(cfg config.Config) flags {
	var opts flags

	flag.StringVar(&opts.modeName, "mode", cfg.Mode, "`MODE` to drill: decode (code to character) or encode (character to code)")
	flag.StringVar(&opts.setName, "set", cfg.Set, "`SET` of characters to drill: letters, digits, alphanumeric or punctuation")
	flag.StringVar(&opts.chars, "chars", cfg.Chars, "drill only these `CHARACTERS`, overriding --set")
	flag.BoolVar(&opts.cheatsheet, "cheatsheet", cfg.Cheatsheet, "show every code of the set under the prompt")
	flag.BoolVar(&opts.debug, "debug", false, "enable debug mode")
	flag.StringVar(&opts.debugFile, "debug-file", logFile, "path to debug log file (ignored if --debug is not set)")
	flag.BoolVar(&opts.showVersion, "version", false, "print the version and exit")
	flag.Parse()

	// A preset picked on the command line beats custom characters from the config.
	if isSet("set") && !isSet("chars") {
		opts.chars = ""
	}

	return opts
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
