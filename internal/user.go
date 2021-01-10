package internal

type User struct {
	Email     string
	Username  string
	Name      string
	AvatarURL string
}

type GithubOAuth struct{}

type SessionOrigin string

const (
	SessionOriginEmpty  = ""
	SessionOriginGithub = "github"
)

type Session struct {
	AccessToken  string
	Email        string
	RefreshToken string
	Origin       SessionOrigin
}
