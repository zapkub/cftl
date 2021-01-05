// +build !prod

package fsutil

import (
	"fmt"
	"log"
	"os"
	"path"
)

func initialdefault() *Manager {
	wd, _ := os.Getwd()

	m := &Manager{
		ExecDir: path.Join(wd, "bin"),
	}
	// checking if the development running in the right place
	// which is workspaceroot
	f, err := os.Open(path.Join(wd, "go.mod"))
	if err != nil {
		log.Fatal(fmt.Errorf("cannot access go.mod. make sure you are running development from root workspace: %w", err))
	}

	f.Close()

	return m
}
