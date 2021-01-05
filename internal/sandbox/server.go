package sandbox

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/zapkub/cftl/internal/fsutil"
	"github.com/zapkub/cftl/internal/logger"
)

type Server struct{}

func (s *Server) Install(m *http.ServeMux) {

	m.Handle("/apis/execute", http.HandlerFunc(s.handleExecute))

}

func (s *Server) handleExecute(w http.ResponseWriter, r *http.Request) {
	log.Println("execute input", r.FormValue("source"))
	cmd := fsutil.Default.Exec("node", "-e", r.FormValue("source"))
	var result, err = cmd.CombinedOutput()
	if err != nil {
		logger.Errorf(r.Context(), "run error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("stdout", string(result))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, bytes.NewBuffer(result))
}
