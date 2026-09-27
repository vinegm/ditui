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

func main() {
	// Loaded first so the flags default to it and only override what they set.
	cfg, cfgErr := config.Load(config.DefaultPath())

	modeName := flag.String("mode", cfg.Mode, "`MODE` to drill: decode (code to character) or encode (character to code)")
	setName := flag.String("set", cfg.Set, "`SET` of characters to drill: letters, digits, alphanumeric or punctuation")
	chars := flag.String("chars", cfg.Chars, "drill only these `CHARACTERS`, overriding --set")
	cheatsheet := flag.Bool("cheatsheet", cfg.Cheatsheet, "show every code of the set under the prompt")
	debug := flag.Bool("debug", false, "enable debug mode")
	debugFile := flag.String("debug-file", logFile, "path to debug log file (ignored if --debug is not set)")
	showVersion := flag.Bool("version", false, "print the version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Println("dit version", version)
		return
	}

	log.SetOutput(io.Discard)
	if *debug {
		file, err := tea.LogToFile(*debugFile, logLevel)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	if cfgErr != nil && !os.IsNotExist(cfgErr) {
		log.Printf("Could not load config: %v", cfgErr)
	}

	// A preset picked on the command line beats custom characters from the config.
	if isSet("set") && !isSet("chars") {
		*chars = ""
	}

	mode, err := morse.ParseMode(*modeName)
	if err != nil {
		usageError(err.Error())
	}

	set, err := morse.ResolveSet(*setName, *chars)
	if err != nil {
		usageError(err.Error())
	}

	options := drill.Options{
		Mode:       mode,
		Set:        set,
		Cheatsheet: *cheatsheet,
	}

	program := tea.NewProgram(tui.New(cfg, options))
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
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
