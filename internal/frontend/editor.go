package frontend

import "net/http"

type editorPage struct{}

func (s *Server) editorPageHandler(w http.ResponseWriter, r *http.Request) {
	s.servePage(r.Context(), w, "editor.html", editorPage{})
}
