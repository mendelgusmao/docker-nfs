package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mendelgusmao/docker-nfs/lib/config"
	"github.com/willscott/memphis"

	nfs "github.com/willscott/go-nfs"
	nfshelper "github.com/willscott/go-nfs/helpers"
)

type Server struct {
	nfs.Server
	listener net.Listener
}

func Create(config config.Config, path string) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:0", config.Iface))

	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v\n", err)
	}

	fs := memphis.FromOS(path)
	bfs := fs.AsBillyFS(0, 0)
	handler := nfshelper.NewNullAuthHandler(bfs)
	cacheHelper := nfshelper.NewCachingHandler(handler, 1024)

	s := &Server{
		listener: listener,
	}
	s.Handler = cacheHelper
	s.Context = context.Background()

	log.Printf("serving %s at %s\n", path, listener.Addr())

	return s, nil
}

func (s *Server) Serve() error {
	return s.Server.Serve(s.listener)
}

func (s *Server) Stop() {
	s.Context.Done()
}

func (s *Server) Address() string {
	return s.listener.Addr().String()
}
