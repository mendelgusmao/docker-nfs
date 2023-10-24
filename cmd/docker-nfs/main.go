package main

import (
	"log"
	"os"

	"github.com/mendelgusmao/docker-nfs/lib/wrapper"
)

func init() {
	log.SetPrefix("[docker-nfs] ")
}

func main() {
	w := wrapper.New(os.Args)
	err := w.Wrap()
	log.Fatalln(err)
}
