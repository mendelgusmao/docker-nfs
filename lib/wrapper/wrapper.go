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
				return fmt.Errorf("failed to get absolute path for `%s`: %v\n", path, err)
			}

			w.orchestrator.CreateServer(absPath)
		}
	}

	volumes, _, err := w.extractVolumes()

	if err != nil {
		return err
	}

	return w.createNFSServersFromVolumes(fixedPaths, volumes)
}

func (w *Wrapper) extractVolumes() ([]*Volume, []string, error) {
	volumes := make([]*Volume, 0)
	skip := make(map[int]any, 0)
	filteredArgs := make([]string, 0)
	copy(filteredArgs, w.args)

	for index, arg := range w.args {
		next := index + 1
		_, isVOption := vOption[arg]

		if isVOption && next < len(w.args) {
			clVolume := w.args[next]
			volume, err := volumeFromVOption(clVolume)

			if err != nil {
				return nil, nil, err
			}

			volumes = append(volumes, volume)
			skip[index] = nil
			skip[next] = nil
		}
	}

	for index, arg := range w.args {
		if _, ok := skip[index]; ok {
			continue
		}

		filteredArgs = append(filteredArgs, arg)
	}

	return volumes, filteredArgs, nil
}

func (w *Wrapper) createNFSServersFromVolumes(fixedPaths []string, volumes []*Volume) error {
	for _, volume := range volumes {
		hasFixedPath := false

		for _, fixedPath := range fixedPaths {
			if strings.HasPrefix(volume.Destination, fixedPath) {
				hasFixedPath = true
				break
			}
		}

		if !hasFixedPath {
			w.orchestrator.CreateServer(volume.Destination)
		}
	}

	return nil
}
