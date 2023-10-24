package orchestrator

import (
	"fmt"

	"github.com/mendelgusmao/docker-nfs/lib/config"
	"github.com/mendelgusmao/docker-nfs/lib/server"
	nfsserver "github.com/mendelgusmao/docker-nfs/lib/server"
)

type Orchestrator struct {
	config  config.Config
	servers map[string]*server.Server
	done    chan any
}

func New(config config.Config) *Orchestrator {
	return &Orchestrator{
		config:  config,
		servers: make(map[string]*server.Server, 0),
		done:    make(chan any, 0),
	}
}

func (o *Orchestrator) CreateServer(path string) (*server.Server, error) {
	if server, ok := o.servers[path]; ok {
		return server, nil
	}

	server, err := nfsserver.Create(o.config, path)

	if err != nil {
		return nil, err
	}

	o.servers[path] = server
	go server.Serve()

	return server, nil
}

func (o *Orchestrator) DestroyServer(path string) error {
	if _, ok := o.servers[path]; !ok {
		return fmt.Errorf("server for folder %s doesnt exist", path)
	}

	o.servers[path].Stop()
	delete(o.servers, path)

	if len(o.servers) == 0 {
		o.done <- nil
	}

	return nil
}

func (o *Orchestrator) Wait() {
	if len(o.servers) == 0 {
		return
	}

	select {
	case <-o.done:
		return
	}
}
