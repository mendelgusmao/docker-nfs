package wrapper

import (
	"fmt"
	"strings"
)

type Volume struct {
	Source      string
	Destination string
	Options     string
}

func volumeFromVOption(arg string) (*Volume, error) {
	parts := strings.Split(arg, ":")

	if len(parts) == 1 {
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

	return &Volume{
		Source:      source,
		Destination: destination,
		Options:     options,
	}, nil
}
