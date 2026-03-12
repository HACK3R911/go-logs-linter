package config

import (
	"errors"
	"fmt"

	"github.com/HACK3R911/go-logs-linter/pkg/analyzer"
	"github.com/spf13/viper"
)

var ErrConfigNotFound = errors.New("config file not found")

type Config struct {
	Rules RulesConfig `yaml:"rules"`
}

type RulesConfig struct {
	AllowUppercaseStart bool     `yaml:"allow_uppercase_start"`
	AllowedPatterns     []string `yaml:"allowed_patterns"`
	DisallowedPatterns  []string `yaml:"disallowed_patterns"`
	AllowNonEnglish     bool     `yaml:"allow_non_english"`
	AllowSpecialChars   bool     `yaml:"allow_special_chars"`
	AllowSensitiveData  bool     `yaml:"allow_sensitive_data"`
	SensitiveKeywords   []string `yaml:"sensitive_keywords"`
}

func Load(path string) (*Config, error) {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("%w: config.yaml not found", ErrConfigNotFound)
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) ToSettings() analyzer.Settings {
	return analyzer.Settings{
		AllowUppercaseStart: c.Rules.AllowUppercaseStart,
		AllowedPatterns:     c.Rules.AllowedPatterns,
		DisallowedPatterns:  c.Rules.DisallowedPatterns,
		AllowNonEnglish:     c.Rules.AllowNonEnglish,
		AllowSpecialChars:   c.Rules.AllowSpecialChars,
		AllowSensitiveData:  c.Rules.AllowSensitiveData,
		SensitiveKeywords:   c.Rules.SensitiveKeywords,
	}
}
