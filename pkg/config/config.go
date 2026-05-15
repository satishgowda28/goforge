package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Server struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}
type Database struct {
	URL string `mapstructure:"db_url"`
}
type LLM struct {
	Provider       string `mapstructure:"provider"`
	Model          string `mapstructure:"model"`
	APIKey         string `mapstructure:"api_key"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}
type Tools struct {
	WorkingDir           string   `mapstructure:"working_dir"`
	ShellAllowedCommands []string `mapstructure:"shell_allowed_commands"`
}
type Worker struct {
	PoolSize int `mapstructure:"pool_size"`
}

type Config struct {
	Server Server   `mapstructure:"server"`
	DB     Database `mapstructure:"db"`
	LLM    LLM      `mapstructure:"llm"`
	Tools  Tools    `mapstructure:"tools"`
	Worker Worker   `mapstructure:"worker"`
}

var cfg *Config

func Load() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file present")
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Options 1
	/* viper.SetEnvPrefix("GOFORGE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() */
	// Options 2
	_ = viper.BindEnv("llm.api_key", "GOFORGE_LLM_API_KEY")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("llm.timeout_seconds", 30)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var loaded Config
	if err := viper.Unmarshal(&loaded); err != nil {
		return nil, fmt.Errorf("unable to decode struct: %w", err)
	}
	cfg = &loaded

	return cfg, nil
}

func Get() *Config {
	if cfg == nil {
		panic("config not loaded — call Load() first")
	}
	return cfg
}
