package frontend

import "net/http"

type indexPage struct{}

func (s *Server) indexPageHandler(w http.ResponseWriter, r *http.Request) {
	s.servePage(r.Context(), w, "index.html", indexPage{})
}
