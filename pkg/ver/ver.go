package ver

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

func Exists(version string) bool {
	home, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln("Could not retrieve user home directory", err)
	}

	_, err = os.Stat(filepath.Join(home, ".nvm/versions/node", version))

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func Get() string {
	home, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln("Could not retrieve user home directory", err)
	}

	currentFile, err := os.Open(filepath.Join(home, ".nvm/current.txt"))

	if err != nil {
		log.Fatalln("Could not open file 'current.txt'", err)
	}

	defer currentFile.Close()

	scanner := bufio.NewScanner(currentFile)

	scanner.Scan()

	current := scanner.Text()

	return current
}
