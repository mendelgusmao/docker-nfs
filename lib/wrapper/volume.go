package wrapper

import (
	"fmt"
	"os"
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

func (v Volume) ToCLOptions() string {
	return fmt.Sprintf(clTemplate, v.ServerAddress, v.Destination, v.Name)
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

	name := strings.Replace(strings.ToLower(source), string(os.PathSeparator), "_", -1)

	return &Volume{
		Source:      source,
		Destination: destination,
		Options:     options,
		Name:        name,
	}, nil
}
