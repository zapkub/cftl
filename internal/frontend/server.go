package frontend

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/safehtml/template"
	"github.com/zapkub/cftl/internal/auth"
	"github.com/zapkub/cftl/internal/fsutil"
	"github.com/zapkub/cftl/internal/logger"
)

type Server struct {
	templateDir   template.TrustedSource
	templates     map[string]*template.Template
	authenticator *auth.Authenticator
}

func New(authenticator *auth.Authenticator) *Server {
	return &Server{
		authenticator: authenticator,
	}
}

func (s *Server) Install(m *http.ServeMux) {
	s.templateDir = template.TrustedSourceJoin(
		fsutil.Default.WebDir(),
		template.TrustedSourceFromConstant("html"),
	)
	s.templates = parsePageTemplates(s.templateDir)

	m.Handle("/editor", http.HandlerFunc(s.editorPageHandler))
	m.Handle("/auth/github_callback", http.HandlerFunc(s.githubCallbackhandler))
	m.Handle("/auth", http.HandlerFunc(s.authPageHandler))
	m.Handle("/", http.HandlerFunc(s.indexPageHandler))

}

func (s *Server) renderPage(ctx context.Context, templateName string, page interface{}) ([]byte, error) {
	tmpl, err := s.findTemplate(templateName)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Server) servePage(ctx context.Context, w http.ResponseWriter, templateName string, page interface{}) {

	buf, err := s.renderPage(ctx, templateName, page)
	if err != nil {
		logger.Errorf(ctx, "render page %q, %+v: %v", templateName, page, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if _, err := io.Copy(w, bytes.NewReader(buf)); err != nil {
		logger.Errorf(ctx, "cannot write template to ResponseWriter: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (s *Server) findTemplate(tempateName string) (*template.Template, error) {
	tmpl := s.templates[tempateName]
	if tmpl == nil {
		return nil, fmt.Errorf("BUG: template not found (%q)", tempateName)
	}

	return tmpl, nil
}

func parsePageTemplates(dir template.TrustedSource) map[string]*template.Template {
	var ts = make(map[string]*template.Template)
	tsc := template.TrustedSourceFromConstant
	join := template.TrustedSourceJoin

	htmlSets := [][]template.TrustedSource{
		{tsc("editor.html")},
		{tsc("auth.html")},
		{tsc("auth_callback.html")},
		{tsc("index.html")},
	}

	for _, set := range htmlSets {
		t, err := template.New("base.html").
			ParseFilesFromTrustedSources(join(dir, tsc("base.html")))
		if err != nil {
			log.Fatalf("cannot parse files: %v", err)
		}

		var files []template.TrustedSource
		for _, f := range set {
			files = append(files, join(dir, tsc("pages"), f))
		}

		if _, err := t.ParseFilesFromTrustedSources(files...); err != nil {
			log.Fatalf("cannot parse files from (%v): %v", files, err)
		}
		ts[set[0].String()] = t
		log.Printf("load html template: %v", set[0].String())
	}

	return ts
}
