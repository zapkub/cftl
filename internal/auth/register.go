package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/zapkub/cftl/internal"
	"github.com/zapkub/cftl/internal/apperror"
	"github.com/zapkub/cftl/internal/logger"
)

type UserInfo interface {
	UserInfo() (Email, Username, Name, DisplayImageURL string)
}

func registerNewUser(ctx context.Context) {

}

func (auth *Authenticator) LoginWithOAuthOrigin(ctx context.Context, origin internal.SessionOrigin, accessToken, refreshToken string, expiresIn int) error {

	var userinfo *internal.User

	switch origin {
	case internal.SessionOriginGithub:
		githubUserInfo, err := retrieveUserInformationFromGithub(ctx, accessToken)
		if err != nil {
			return err
		}
		userinfo = githubUserInfo.UserInfo()
	}

	// will automatically register if user does not have an account
	// but in future we should have a saperate flow for register
	if _, err := auth.db.GetUser(ctx, userinfo.Email); err != nil && errors.Is(err, apperror.NotFound) {
		_, err = auth.db.InsertUser(ctx, userinfo)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err := auth.db.InsertSession(ctx, &internal.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        userinfo.Email,
		Origin:       origin,
	})
	if err != nil {
		return err
	}

	return nil
}

const githubGetUserEndpoint = "https://api.github.com/user"

type githubUserResponseBody struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

func (g *githubUserResponseBody) UserInfo() *internal.User {
	return &internal.User{
		Email:     g.Email,
		Username:  g.Login,
		Name:      g.Name,
		AvatarURL: g.AvatarURL,
	}
}

func retrieveUserInformationFromGithub(ctx context.Context, accessToken string) (*githubUserResponseBody, error) {

	var req, err = http.NewRequest("GET", githubGetUserEndpoint, nil)
	if err != nil {
		logger.Errorf(ctx, "BUG: create request for github user api error: %+v", err)
		return nil, fmt.Errorf("cannot create new request for github user info: %w", err)
	}
	req.Header.Set("authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf(ctx, "github user API return error (%d): %+v", resp.StatusCode, err)
		return nil, fmt.Errorf("cannot get Github user info: %w", err)
	}
	var respBody githubUserResponseBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		logger.Errorf(ctx, "cannot parse response from github user api: %+v", err)
		return nil, fmt.Errorf("parse result from github user api error:%w", err)
	}

	return &respBody, nil
}
