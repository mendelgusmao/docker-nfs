package wrapper

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mendelgusmao/docker-nfs/lib/orchestrator"
)

var (
	vOption = map[string]any{
		"-v":       nil,
		"--volume": nil,
	}
)

type Wrapper struct {
	args         []string
	orchestrator *orchestrator.Orchestrator
}

func New(args []string) *Wrapper {
	return &Wrapper{
		args:         args,
		orchestrator: orchestrator.New(),
	}
}

func (w *Wrapper) Wrap() error {
	if err := w.createNFSServers(); err != nil {
		return err
	}

	w.orchestrator.Wait()

	return nil
}

func (w *Wrapper) createNFSServers() error {
	fixedPaths, ok := tryLoadingDockerNFSFile()

	if ok {
		for _, path := range fixedPaths {
			absPath, err := filepath.Abs(path)

			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s: %v\n", path, err)
			}

			w.orchestrator.CreateServer(absPath)
		}
	}

	for index, arg := range w.args {
		next := index + 1
		_, isVOption := vOption[arg]

		if isVOption && next < len(w.args) {
			clVolume := w.args[next]
			err := w.createNFSServerFromVOption(fixedPaths, clVolume)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Wrapper) createNFSServerFromVOption(fixedPaths []string, clVolume string) error {
	hasFixedPath := false
	volume, err := volumeFromVOption(clVolume)

	if err != nil {
		return err
	}

	for _, fixedPath := range fixedPaths {
		if strings.HasPrefix(volume.Destination, fixedPath) {
			hasFixedPath = true
			continue
		}
	}

	if !hasFixedPath {
		w.orchestrator.CreateServer(volume.Destination)
	}

	return nil
}
