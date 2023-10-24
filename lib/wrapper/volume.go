package wrapper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const clTemplate = "--opt type=nfs --opt o=addr=%s,rw --opt device=:%s %s"

type Volume struct {
	Source        string
	Destination   string
	Options       string
	ServerAddress string
	Name          string
}

func (v Volume) ToCLOptions() []string {
	return []string{
		"--opt",
		"type=nfs",
		"--opt",
		fmt.Sprintf("o=addr=%s,rw", v.ServerAddress),
		"--opt",
		fmt.Sprintf("device=:%s", v.Destination),
		v.Name,
	}
}

func volumeFromVOption(arg string) (*Volume, error) {
	parts := strings.Split(arg, ":")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid volume: `%s`", arg)
	}

	var (
		source      string
		destination string
		options     string
	)

	if len(parts) >= 2 {
		source = parts[0]
		destination = parts[1]
	}

	if len(parts) == 3 {
		options = parts[2]
	}

	absSource, err := filepath.Abs(source)

	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for `%s`: %v\n", source, err)
	}

	name := strings.Trim(
		strings.ReplaceAll(
			strings.ToLower(absSource),
			string(os.PathSeparator),
			"_",
		),
		"_",
	)

	return &Volume{
		Source:      absSource,
		Destination: destination,
		Options:     options,
		Name:        name,
	}, nil
}
