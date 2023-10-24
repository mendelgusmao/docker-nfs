package main

import (
	"log"
	"os"

	"github.com/mendelgusmao/docker-nfs/lib/config"
	"github.com/mendelgusmao/docker-nfs/lib/wrapper"
)

func init() {
	log.SetPrefix("[docker-nfs] ")
}

func main() {
	config, err := config.Load()

	if err != nil {
		log.Fatalln(err)
	}

	w := wrapper.New(*config, os.Args)
	err = w.Wrap()
	log.Fatalln(err)
}
