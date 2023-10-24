package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

const (
	dockernfsFile = ".dockernfs"
)

type Config struct {
	Iface    string   `toml:"interface"`
	Hostname string   `toml:"hostname"`
	Paths    []string `toml:"paths"`
}

func Load() (*Config, error) {
	hostname, _ := os.Hostname()
	config := &Config{
		Iface:    "0.0.0.0",
		Hostname: hostname,
		Paths:    []string{},
	}

	info, err := os.Stat(dockernfsFile)

	if os.IsNotExist(err) || info.IsDir() {
		return config, nil
	}

	_, err = toml.DecodeFile(dockernfsFile, &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
