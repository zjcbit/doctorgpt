package main

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

// TODO: Make error placeholder configurable
// TODO: Send which image, program and/or version is outputing the logs (if known)
const errorPlaceholder = "$ERROR"

var basePrompt = "You are ErrorDebuggingGPT. Your sole purpose in this world is to help software engineers by diagnosing software system errors and bugs that can occur in any type of computer system. The message following the first line containing \"ERROR:\" up until the end of the prompt is a computer error no more and no less. It is your job to try to diagnose and fix what went wrong. Ready?\nERROR:\n" + errorPlaceholder

type config struct {
	Prompt  string         `yaml:"prompt,omitempty"`
	Parsers []parserConfig `yaml:"parsers"`
}

type parserConfig struct {
	Regex    string            `yaml:"regex"`
	Triggers []variableMatcher `yaml:"triggers,omitempty"`
	Filters  []variableMatcher `yaml:"filters,omitempty"`
	Excludes []variableMatcher `yaml:"excludes,omitempty"`
}

type variableMatcher struct {
	Variable string `yaml:"variable"`
	Regex    string `yaml:"regex"`
}

type configProvider func(log *zap.SugaredLogger, configFile string) (config, error)

func fileConfigProvider(log *zap.SugaredLogger, configFile string) (config, error) {
	// Read configuration
	var config config
	bytes, err := readBytes(configFile)
	if err != nil {
		return config, fmt.Errorf("Failed to open config file: %w", err)
	}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return config, fmt.Errorf("Invalid config: %w", err)
	}
	return config, nil
}

func readBytes(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
