package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// OAuthConfig defines enviroment variables for OAuth
type OAuthConfig struct {
	ClientID     string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"CLIENT_SECRET" required:"true"`
}

// OAuth service abstracts OAuth implementation
type OAuth struct {
	config *oauth2.Config
	store  *sessions.CookieStore
}

// NewOAuth return new OAuth service
func NewOAuth(clientID, clientSecret string) *OAuth {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"read:user"},
		Endpoint:     github.Endpoint,
	}
	return &OAuth{
		config: config,
		store:  sessions.NewCookieStore([]byte(clientSecret)),
	}
}

// GithubUser represents the user response returned by the GitHub auth.
type githubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Username  string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// MakeAuthURL returns string for redirect to provider
func (o *OAuth) MakeAuthURL(w http.ResponseWriter, r *http.Request) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	session, _ := o.store.Get(r, "sess")
	session.Values["state"] = state
	session.Save(r, w)

	return o.config.AuthCodeURL(state)
}

// ValidateState protects the user from CSRF attacks
func (o *OAuth) ValidateState(r *http.Request, state string) error {
	session, err := o.store.Get(r, "sess")
	if err != nil {
		return fmt.Errorf("can't get session: %s", err)
	}
	if state != session.Values["state"] {
		return fmt.Errorf("incorrect state: %s", state)
	}
	return nil
}

// GetUser gets user from provider and return user model
func (o *OAuth) GetUser(ctx context.Context, code string) (*githubUser, error) {
	token, err := o.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth exchange error: %s", err)
	}
	client := o.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("can't get user from github: %s", err)
	}
	defer resp.Body.Close()
	var user githubUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("can't parse github response: %s", err)
	}

	return &user, nil
}
