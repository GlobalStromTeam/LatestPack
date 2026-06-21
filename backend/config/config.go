package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Port           string        `yaml:"port"`
	JWTSecret      string        `yaml:"jwt_secret"`
	JWTTTL         time.Duration `yaml:"-"`
	JWTTTLH        int           `yaml:"jwt_ttl_hours"`
	DataDir        string        `yaml:"data_dir"`
	AllowedOrigins []string      `yaml:"allowed_origins"`
}

func LoadConfig() *Config {
	cfg := defaults()

	configPath := filepath.Join(cfg.DataDir, "config.yaml")
	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: parse config.yaml: %v\n", err)
			cfg = defaults()
		}
	} else {
		if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
			log.Printf("Warning: create data dir: %v", err)
		}
		writeDefaultConfig(cfg.DataDir, cfg)
	}

	cfg.JWTTTL = time.Duration(cfg.JWTTTLH) * time.Hour

	if len(cfg.AllowedOrigins) == 0 {
		cfg.AllowedOrigins = []string{"*"}
	}

	return cfg
}

func generateSecret() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("Failed to generate secret: %v", err)
	}
	return hex.EncodeToString(b)
}

func defaults() *Config {
	return &Config{
		Port:           ":8080",
		JWTSecret:      generateSecret(),
		JWTTTLH:        6,
		DataDir:        "./data",
		AllowedOrigins: []string{"*"},
	}
}

func writeDefaultConfig(dataDir string, cfg *Config) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return
	}
	header := "# LatestPack configuration\n"
	if err := os.WriteFile(filepath.Join(dataDir, "config.yaml"), append([]byte(header), data...), 0600); err != nil {
		log.Printf("Warning: write default config: %v", err)
	}
}
