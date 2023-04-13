package config

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/hcl/v2"
	"github.com/go-playground/mold/v4/modifiers"
)

type Config struct {
	Actor     string           `hcl:"actor" mod:"trim"`
	Output    OutputConfig     `hcl:"output,block"`
	Resources []ResourceConfig `hcl:"resource,block" mod:"dive"`
}

type OutputConfig struct {
	Path    string `hcl:"path" validate:"required" mod:"trim"`
	Package string `hcl:"package" validate:"required" mod:"trim"`
}

type ResourceConfig struct {
	Name        string            `hcl:"name,label" mod:"trim,snake"`
	Model       string            `hcl:"model" mod:"trim"`
	Permissions []string          `hcl:"permissions" mod:"dive,trim,snake"`
	Roles       []RoleConfig      `hcl:"role,block" mod:"dive"`
	Attributes  []AttributeConfig `hcl:"attribute,block" mod:"dive"`
}

type RoleConfig struct {
	Name        string   `hcl:"name,label" mod:"trim,snake"`
	Permissions []string `hcl:"permissions" mod:"dive,trim,snake"`
}

type AttributeConfig struct {
	Name        string   `hcl:"name,label" mod:"trim,snake"`
	Permissions []string `hcl:"permissions" mod:"dive,trim,snake"`
}

func ReadConfig(data []byte) (*Config, error) {
	var config Config

	err := hcl.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Use mold to transform the config
	m := modifiers.New()
	err = m.Struct(context.Background(), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to transform config: %w", err)
	}

	return &config, nil
}

func LoadConfig(path string) (*Config, error) {
	// Load the config from file
	// TODO: Allow configurable config file path / maybe stdin?
	// err := hclsimple.Decode("", file, nil, &config)

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	return ReadConfig(b)
}
