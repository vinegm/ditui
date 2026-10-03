package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/vinegm/ditui/internal/morse"
	"github.com/vinegm/ditui/internal/storage/paths"
)

// validate reports an error for any setting
// the rest of the program can't work with.
func (config Config) validate() error {
	if _, err := morse.ParseMode(config.Mode); err != nil {
		return fmt.Errorf("mode: %w", err)
	}

	if _, err := morse.ResolveSet(config.Set, config.Chars); err != nil {
		return fmt.Errorf("set/chars: %w", err)
	}

	return nil
}

// Loads configuration from file. The returned config is always usable:
// on error the defaults are returned alongside it.
func Load(path string) (Config, error) {
	config := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Default(), err
	}

	if err := config.validate(); err != nil {
		return Default(), err
	}

	return config, nil
}

// DefaultPath is the config file's path.
func DefaultPath() string {
	return paths.Config("config.yaml")
}
