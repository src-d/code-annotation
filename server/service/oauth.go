package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/src-d/code-annotation/server/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// State is a token to protect the user from CSRF attacks.
// You must always provide a non-empty string and validate that it matches the the state query parameter on your redirect callback
// FIXME issue: https://github.com/src-d/code-annotation/issues/21
const oauthStateString = "state"

// OAuthConfig defines enviroment variables for OAuth
type OAuthConfig struct {
	ClientID     string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"CLIENT_SECRET" required:"true"`
}

// OAuth service abstracts OAuth implementation
type OAuth struct {
	config *oauth2.Config
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
	}
}

type githubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Username  string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// MakeAuthURL returns string for redirect to provider
func (o *OAuth) MakeAuthURL() string {
	return o.config.AuthCodeURL(oauthStateString)
}

// ValidateState protects the user from CSRF attacks
func (o *OAuth) ValidateState(state string) error {
	if state != oauthStateString {
		return fmt.Errorf("incorrect state: %s", state)
	}
	return nil
}

// GetUser gets user from provider and return user model
func (o *OAuth) GetUser(ctx context.Context, code string) (*model.User, error) {
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
	return &model.User{
		ID:        user.ID,
		Login:     user.Login,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
	}, nil
}
