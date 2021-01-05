package fsutil

import (
	"os"
	"os/exec"
	"path"

	"github.com/google/safehtml/template"
)

type Manager struct {
	ExecDir string
}

const defaultprem os.FileMode = 0755

func (m *Manager) OpenResource(name string) (*os.File, error) {
	return os.OpenFile(path.Join(m.ExecDir, name), os.O_RDONLY, defaultprem)
}

func (m *Manager) WebDir() template.TrustedSource {
	_, ok := os.LookupEnv("WEB_DIR")
	if !ok {
		return template.TrustedSourceFromConstant("web")
	}
	return template.TrustedSourceFromEnvVar("WEB_DIR")
}

func (m *Manager) Exec(name string, args ...string) *exec.Cmd {
	return exec.Command(path.Join(m.ExecDir, name), args...)
}

var Default = initialdefault()
