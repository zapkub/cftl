// +build !prod

package fsutil

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/zapkub/cftl/internal/logger"
)

func initialdefault() *Manager {
	wd, _ := os.Getwd()
	var appdir = path.Join(wd, ".cftl")
	err := os.MkdirAll(appdir, Defaultprem)
	if err != nil {
		logger.Fatalf(nil, "cannot initlize app dir (%s): %v", appdir, err)
	}

	m := &Manager{
		ExecDir: path.Join(wd, "bin"),

		// appdir is $rootWorkspace/.cftl
		AppDir: appdir,
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
