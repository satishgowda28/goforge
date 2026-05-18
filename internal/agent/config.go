package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type AgentConfig struct {
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	SystemPrompt string   `yaml:"system_prompt"`
	Tools        []string `yaml:"tools"`
	MaxSteps     int      `yaml:"max_steps"`
	MemoryType   string   `yaml:"memory_type"`
}

type AgentRegistry struct {
	configs map[string]AgentConfig
}

func NewAgentRegistry(dir string) (*AgentRegistry, error) {
	pattern := filepath.Join(dir, "*.yaml")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("Error in reading the file glob %w", err)
	}
	agentConfigs := make(map[string]AgentConfig)

	for _, path := range paths {
		config, err := loadOnceConfig(path)
		if err != nil {
			return nil, err
		}
		agentConfigs[config.Name] = config
	}
	return &AgentRegistry{configs: agentConfigs}, nil
}

func loadOnceConfig(path string) (AgentConfig, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return AgentConfig{}, fmt.Errorf("error in reading file %w", err)
	}
	var config AgentConfig
	if err := yaml.Unmarshal(fileData, &config); err != nil {
		return AgentConfig{}, fmt.Errorf("error in decoding file %w", err)
	}

	return config, nil

}

func (r *AgentRegistry) Get(name string) (AgentConfig, error) {
	config, ok := r.configs[name]
	if !ok {
		return AgentConfig{}, fmt.Errorf("no config for agent %s", name)
	}
	return config, nil
}
