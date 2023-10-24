package wrapper

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/mendelgusmao/docker-nfs/lib/config"
	"github.com/mendelgusmao/docker-nfs/lib/orchestrator"
)

var (
	vOption = map[string]any{
		"-v":       nil,
		"--volume": nil,
	}
	dockerVolumeCreateArgs = []string{"volume", "create", "--driver", "local"}
)

type Wrapper struct {
	config       config.Config
	args         []string
	orchestrator *orchestrator.Orchestrator
	fixedPaths   []string
	volumes      []*Volume
}

func New(config config.Config, args []string) *Wrapper {
	fixedPaths := make([]string, len(config.Paths))
	copy(fixedPaths, config.Paths)

	return &Wrapper{
		config:       config,
		fixedPaths:   fixedPaths,
		args:         args,
		orchestrator: orchestrator.New(config),
	}
}

func (w *Wrapper) Wrap() error {
	_, err := w.createNFSServers()

	if err != nil {
		return err
	}

	if err = w.createNFSVolumes(); err != nil {
		return err
	}

	w.orchestrator.Wait()
	return nil
}

func (w *Wrapper) createNFSVolumes() error {
	for _, volume := range w.volumes {
		log.Printf("creating volume `%s` for `%s`\n", volume.Name, volume.Destination)

		args := append(dockerVolumeCreateArgs, volume.ToCLOptions()...)
		err := execCommand(w.config.Docker.CLIPath, args...)

		if err != nil {
			return err
		}

		log.Printf("volume `%s` created\n", volume.Name)
	}

	return nil
}

func (w *Wrapper) createNFSServers() ([]string, error) {
	for index, path := range w.fixedPaths {
		absPath, err := filepath.Abs(path)

		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path for `%s`: %v\n", path, err)
		}

		w.fixedPaths[index] = absPath
	}

	filteredArgs, err := w.extractVolumes()

	if err != nil {
		return nil, err
	}

	err = w.createNFSServersFromVolumes()

	if err != nil {
		return nil, err
	}

	for _, v := range w.volumes {
		log.Printf("%+#v\n", v)
	}

	return filteredArgs, err
}

func (w *Wrapper) extractVolumes() ([]string, error) {
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
				return nil, err
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

	w.volumes = volumes

	return filteredArgs, nil
}

func (w *Wrapper) createNFSServersFromVolumes() error {
	for _, volume := range w.volumes {
		serverPath := volume.Source

		for _, fixedPath := range w.fixedPaths {
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
