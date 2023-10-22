package orchestrator

import (
	"fmt"

	"github.com/mendelgusmao/docker-nfs/lib/server"
	nfsserver "github.com/mendelgusmao/docker-nfs/lib/server"
)

type Orchestrator struct {
	servers map[string]*server.Server
	done    chan any
}

func New() *Orchestrator {
	return &Orchestrator{
		servers: make(map[string]*server.Server, 0),
		done:    make(chan any, 0),
	}
}

func (o *Orchestrator) CreateServer(path string) error {
	if _, ok := o.servers[path]; ok {
		return fmt.Errorf("server for %s is already registered", path)
	}

	server, err := nfsserver.Create(path)

	if err != nil {
		return err
	}

	o.servers[path] = server
	go server.Serve()

	return nil
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