package wrapper

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const (
	dockernfsFile = ".dockernfs"
)

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
