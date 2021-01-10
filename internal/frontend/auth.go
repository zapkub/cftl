package frontend

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zapkub/cftl/internal"
	"github.com/zapkub/cftl/internal/conf"
	"github.com/zapkub/cftl/internal/logger"
)

type authPage struct {
	GithubClientID string
}

func (s *Server) authPageHandler(w http.ResponseWriter, r *http.Request) {
	s.servePage(r.Context(), w, "auth.html", authPage{
		GithubClientID: conf.C.Oauth.GithubClientID,
	})
}

const (
	githubOAuthEndpoint = "https://github.com/login/oauth/access_token"
)

func (s *Server) githubCallbackhandler(w http.ResponseWriter, r *http.Request) {
	var logincode = r.FormValue("code")
	data := url.Values{}
	data.Add("client_id", conf.C.Oauth.GithubClientID)
	data.Add("client_secret", conf.C.Oauth.GithubClientSecret)
	data.Add("code", logincode)
	data.Add("accept", "application/json")

	var req, _ = http.NewRequest("POST", githubOAuthEndpoint, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf(r.Context(), "github callback unexpected error: %+v", err)
		return
	}

	type resbody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expire_in"`
		Scope        string `json:"scope"`
	}

	decoder := json.NewDecoder(res.Body)
	var responseBody resbody
	err = decoder.Decode(&responseBody)
	if err != nil {
		logger.Errorf(r.Context(), "parsing response from github auth error: %+v", err)
		return
	}

	// retrieve user information
	// what we need is just email and username
	err = s.authenticator.LoginWithOAuthOrigin(
		r.Context(), internal.SessionOriginGithub,
		responseBody.AccessToken,
		responseBody.RefreshToken,
		responseBody.ExpiresIn,
	)
	if err != nil {
		return
	}

	s.servePage(r.Context(), w, "auth_callback.html")
}
