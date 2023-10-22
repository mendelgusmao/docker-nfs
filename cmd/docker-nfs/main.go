package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mendelgusmao/docker-nfs/lib/orchestrator"
)

const dockernfsFile = ".dockernfs"

func init() {
	log.SetPrefix("[docker-nfs] ")
}

func main() {
	o := orchestrator.New()

	if paths, ok := tryLoadingDockerNFSFile(); ok {
		for _, path := range paths {
			absPath, err := filepath.Abs(path)

			if err != nil {
				log.Printf("failed to get absolute path for %s: %v\n", path, err)
				continue
			}

			o.CreateServer(absPath)
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
