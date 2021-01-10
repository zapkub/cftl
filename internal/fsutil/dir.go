// +build prod

package fsutil

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

func initialdefault() *Manager {
	execpath, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}

	execpath, err = filepath.EvalSymlinks(execpath)
	if err != nil {
		log.Fatal(err)
	}

	return &Manager{
		ExecDir: path.Dir(execpath),
		AppDir:  createAppDir(),
	}
}
