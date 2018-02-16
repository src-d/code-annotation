package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/src-d/code-annotation/server/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// ErrNoAccess means user doesn't have necessary permissions for a resource
var ErrNoAccess = errors.New("access denied")

// OAuthConfig defines enviroment variables for OAuth
type OAuthConfig struct {
	ClientID                string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret            string `envconfig:"CLIENT_SECRET" required:"true"`
	RestrictAccess          string `envconfig:"RESTRICT_ACCESS"`
	RestrictRequesterAccess string `envconfig:"RESTRICT_REQUESTER_ACCESS"`
}

// Needed for the tests
type restrictionChecker interface {
	checkAccess(client *http.Client, restriction, login string) error
}

// OAuth service abstracts OAuth implementation
type OAuth struct {
	config                  *oauth2.Config
	store                   *sessions.CookieStore
	restrictAccess          string
	restrictRequesterAccess string
	restrictionChecker      restrictionChecker
}

// NewOAuth return new OAuth service
func NewOAuth(clientID, clientSecret, restrictAccess, restrictRequesterAccess string) *OAuth {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"read:user", "read:org"},
		Endpoint:     github.Endpoint,
	}
	return &OAuth{
		config:                  config,
		store:                   sessions.NewCookieStore([]byte(clientSecret)),
		restrictAccess:          restrictAccess,
		restrictRequesterAccess: restrictRequesterAccess,
		restrictionChecker:      &githubPermissions{},
	}
}

// GithubUser represents the user response returned by the GitHub auth.
type githubUser struct {
	ID        int        `json:"id"`
	Login     string     `json:"login"`
	Username  string     `json:"name"`
	AvatarURL string     `json:"avatar_url"`
	Role      model.Role `json:"-"`
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
		return nil, fmt.Errorf("can't parse github user response: %s", err)
	}

	if err = o.setRole(client, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// setRole sets the Role attribute of the given GitHub user based on the
// organization and team memberships
func (o *OAuth) setRole(client *http.Client, user *githubUser) error {
	defaultRole := model.Requester

	if o.restrictRequesterAccess != "" {
		if err := o.restrictionChecker.checkAccess(client, o.restrictRequesterAccess, user.Login); err == nil {
			user.Role = model.Requester
			return nil
		} else if err != ErrNoAccess {
			return err
		}

		defaultRole = model.Worker
	}

	if o.restrictAccess != "" {
		if err := o.restrictionChecker.checkAccess(client, o.restrictAccess, user.Login); err != nil {
			return err
		}
	}

	user.Role = defaultRole
	return nil
}

const orgPrefix = "org:"
const teamPrefix = "team:"

// githubPermissions implements the restrictionChecker interface
type githubPermissions struct{}

func (o *githubPermissions) checkAccess(client *http.Client, restriction, login string) error {
	if strings.HasPrefix(restriction, orgPrefix) {
		org := strings.TrimPrefix(restriction, orgPrefix)
		return o.checkUserInOrg(client, org, login)
	}

	if strings.HasPrefix(restriction, teamPrefix) {
		team := strings.TrimPrefix(restriction, teamPrefix)
		return o.checkUserInTeam(client, team, login)
	}

	return fmt.Errorf("invalid restriction '%s', it must be one of 'org:<organization-name>' or 'team:<team-id>'", restriction)
}

func (o *githubPermissions) checkUserInOrg(client *http.Client, org, login string) error {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/members/%s", org, login)
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("can't get user organizations from github: %s", err)
	}

	// StatusNoContent means user is a member
	// https://developer.github.com/v3/orgs/members/#check-membership
	if resp.StatusCode != http.StatusNoContent {
		return ErrNoAccess
	}

	return nil
}

func (o *githubPermissions) checkUserInTeam(client *http.Client, team, login string) error {
	url := fmt.Sprintf("https://api.github.com/teams/%s/memberships/%s", team, login)
	r, err := http.NewRequest("GET", url, nil)
	// The Nested Teams API is currently available for developers to preview only
	r.Header.Add("Accept", "application/vnd.github.hellcat-preview+json")
	if err != nil {
		return fmt.Errorf("can't create team request: %s", err)
	}

	resp, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("can't get user organizations from github: %s", err)
	}

	// Only StatusOK means user is a member
	// https://developer.github.com/v3/orgs/teams/#get-team-membership
	if resp.StatusCode != http.StatusOK {
		return ErrNoAccess
	}

	// Don't allow access to pending members
	defer resp.Body.Close()
	var teamResp struct {
		State string `json:"state"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&teamResp); err != nil {
		return fmt.Errorf("can't parse github team response: %s", err)
	}

	if teamResp.State != "active" {
		return ErrNoAccess
	}

	return nil
}
