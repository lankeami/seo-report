package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type RSSSource struct {
	Name   string `mapstructure:"name"`
	URL    string `mapstructure:"url"`
	Weight int    `mapstructure:"weight"`
}

type SourcesConfig struct {
	RSS []RSSSource `mapstructure:"rss"`
}

// CategoryConfig defines a classification category.
// Using a slice (not map) so Viper does not lowercase the Name field.
type CategoryConfig struct {
	Name     string   `mapstructure:"name"`
	Keywords []string `mapstructure:"keywords"`
	Priority int      `mapstructure:"priority"`
}

type DedupConfig struct {
	LevenshteinThreshold float64 `mapstructure:"levenshtein_threshold"`
}

type OutputConfig struct {
	Dir string `mapstructure:"dir"`
}

type Config struct {
	Sources    SourcesConfig       `mapstructure:"sources"`
	Keywords   map[string][]string `mapstructure:"keywords"`
	Categories []CategoryConfig    `mapstructure:"categories"`
	Dedup      DedupConfig         `mapstructure:"dedup"`
	Output     OutputConfig        `mapstructure:"output"`
}

// CategoriesByName returns a name→CategoryConfig map for O(1) lookup.
func (c *Config) CategoriesByName() map[string]CategoryConfig {
	m := make(map[string]CategoryConfig, len(c.Categories))
	for _, cat := range c.Categories {
		m[cat.Name] = cat
	}
	return m
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	return &cfg, nil
}
