package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Context   map[string]interface{} `yaml:"context"`
	Templates map[string]string
	Env       []string
	Exec      []string
}

// Load the configuration from the provided YAML data.
func (cfg *Config) Load(data []byte) error {
	return yaml.Unmarshal(data, cfg)
}

// Read the configuration from the provided YAML file.
func (cfg *Config) Read(file string) error {
	if data, err := ioutil.ReadFile(file); err == nil {
		return cfg.Load(data)
	} else {
		return err
	}
}
