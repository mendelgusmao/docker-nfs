package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mendelgusmao/docker-nfs/lib/orchestrator"
)

const (
	dockernfsFile = ".dockernfs"
	vOption       = "-v"
)

func init() {
	log.SetPrefix("[docker-nfs] ")
}

func main() {
	o := orchestrator.New()
	fixedPaths, ok := tryLoadingDockerNFSFile()

	if ok {
		for _, path := range fixedPaths {
			absPath, err := filepath.Abs(path)

			if err != nil {
				log.Printf("failed to get absolute path for %s: %v\n", path, err)
				continue
			}

			o.CreateServer(absPath)
		}
	}

	for index, arg := range os.Args {
		next := index + 1

		if arg == vOption && next < len(os.Args) {
			clPath := os.Args[next]
			hasFixedPath := false

			for _, fixedPath := range fixedPaths {
				if strings.HasPrefix(clPath, fixedPath) {
					hasFixedPath = true
					continue
				}
			}

			if !hasFixedPath {
				o.CreateServer(clPath)
			}
		}
	}

	o.Wait()
}

func tryLoadingDockerNFSFile() ([]string, bool) {
	info, err := os.Stat(dockernfsFile)

	if os.IsNotExist(err) || info.IsDir() {
		return nil, false
	}

	readFile, err := os.Open(dockernfsFile)

	if err != nil {
		log.Printf("error loading %s: %v\n", dockernfsFile, err)
		return nil, false
	}

	defer readFile.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(readFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) > 0 {
			lines = append(lines, line)
		}
	}

	return lines, len(lines) > 0
}
