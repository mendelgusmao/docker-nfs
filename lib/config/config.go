package config

import (
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
)

const (
	dockernfsFile = ".dockernfs"
)

type Config struct {
	Iface      string   `toml:"interface"`
	Hostname   string   `toml:"hostname"`
	DockerPath string   `toml:"docker"`
	Paths      []string `toml:"paths"`
}

func Load() (*Config, error) {
	dockerPath, err := exec.LookPath("docker")

	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()
	config := &Config{
		Iface:      "0.0.0.0",
		Hostname:   hostname,
		DockerPath: dockerPath,
		Paths:      []string{},
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
