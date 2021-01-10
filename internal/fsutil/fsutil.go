package fsutil

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/google/safehtml/template"
	"github.com/zapkub/cftl/internal/logger"
)

type Manager struct {
	ExecDir string
	AppDir  string
}

const Defaultprem os.FileMode = 0755

func (m *Manager) OpenResource(name string) (*os.File, error) {
	return os.OpenFile(path.Join(m.ExecDir, "..", name), os.O_RDONLY, Defaultprem)
}
func (m *Manager) MustOpenResource(name string) *os.File {

	f, err := m.OpenResource(name)
	if err != nil {
		panic(fmt.Sprintf("cannot open file %q: %v", name, err))
	}
	return f

}

func (m *Manager) OpenFile(name string, flag int) (*os.File, error) {
	return os.OpenFile(path.Join(m.AppDir, name), flag, Defaultprem)
}

func (m *Manager) MustOpenFile(name string, flag int) *os.File {
	f, err := m.OpenFile(name, flag)
	if err != nil {
		panic(fmt.Sprintf("cannot open file %q: %v", name, err))
	}
	return f
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

func (m *Manager) MigrationsDir() string {
	return path.Join(m.ExecDir, "migrations")
}

var Default = initialdefault()

func createAppDir() string {
	h, _ := os.UserHomeDir()
	var appdir = path.Join(h, ".cftl")
	err := os.MkdirAll(appdir, Defaultprem)
	if err != nil {
		logger.Fatalf(nil, "cannot initlize app dir (%s): %v", appdir, err)
	}
	return appdir
}
