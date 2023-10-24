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
	_, err := w.createNFSServers()

	if err != nil {
		return err
	}

	w.orchestrator.Wait()
	return nil
}

func (w *Wrapper) createNFSServers() ([]string, error) {
	fixedPaths, ok := tryLoadingDockerNFSFile()

	if ok {
		for index, path := range fixedPaths {
			absPath, err := filepath.Abs(path)

			if err != nil {
				return nil, fmt.Errorf("failed to get absolute path for `%s`: %v\n", path, err)
			}

			fixedPaths[index] = absPath
		}
	}

	volumes, filteredArgs, err := w.extractVolumes()

	if err != nil {
		return nil, err
	}

	return filteredArgs, w.createNFSServersFromVolumes(fixedPaths, volumes)
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
		serverPath := volume.Destination

		for _, fixedPath := range fixedPaths {
			if strings.HasPrefix(volume.Destination, fixedPath) {
				serverPath = fixedPath
				break
			}
		}

		server, err := w.orchestrator.CreateServer(serverPath)

		if err != nil {
			return err
		}

		volume.ServerAddress = server.Address()
	}

	return nil
}
