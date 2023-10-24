package config

import (
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
)

const (
	dockernfsFile = ".dockernfs"
)

type Docker struct {
	Host    string `toml:"host"`
	CLIPath string `toml:"clipath"`
}

type Config struct {
	Iface    string   `toml:"interface"`
	Hostname string   `toml:"hostname"`
	Docker   Docker   `toml:"docker"`
	Paths    []string `toml:"paths"`
}

func Load() (*Config, error) {
	dockerCLIPath, err := exec.LookPath("docker")

	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()
	config := &Config{
		Iface:    "0.0.0.0",
		Hostname: hostname,
		Docker: Docker{
			Host:    os.Getenv("DOCKER_HOST"),
			CLIPath: dockerCLIPath,
		},
		Paths: []string{},
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
