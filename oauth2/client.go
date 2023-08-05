package oauth2

import (
	"errors"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
)

var (
	// ErrStateNotFound is returned when the state is not found in the SessionController.
	ErrStateNotFound = errors.New("state could not be found")

	// ErrSessionExpired is returned when the Session has expired.
	ErrSessionExpired = errors.New("access token expired. refresh the session")

	// ErrMissingOAuth2Scope is returned when a specific Rest scope is missing.
	ErrMissingOAuth2Scope = func(scope discord.OAuth2Scope) error {
		return fmt.Errorf("missing '%s' scope", scope)
	}
)

// Session represents a discord access token response (https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-access-token-response)
type Session struct {
	// AccessToken allows requesting user information
	AccessToken string `json:"access_token"`

	// RefreshToken allows refreshing the AccessToken
	RefreshToken string `json:"refresh_token"`

	// Scopes returns the discord.OAuth2Scope(s) of the Session
	Scopes []discord.OAuth2Scope `json:"scopes"`

	// TokenType returns the discord.TokenType of the AccessToken
	TokenType discord.TokenType `json:"token_type"`

	// Expiration returns the time.Time when the AccessToken expires and needs to be refreshed
	Expiration time.Time `json:"expiration"`
}

func (s Session) Expired() bool {
	return s.Expiration.Before(time.Now())
}

// New returns a new Rest client with the given ID, secret and ConfigOpt(s).
func New(id snowflake.ID, secret string, opts ...ConfigOpt) *Client {
	config := DefaultConfig()
	config.Apply(opts)

	return &Client{
		ID:              id,
		Secret:          secret,
		Rest:            config.Rest,
		StateController: config.StateController,
	}
}

// Client is a high level wrapper around Discord's OAuth2 API.
type Client struct {
	ID              snowflake.ID
	Secret          string
	Rest            rest.OAuth2
	StateController StateController
}

// GenerateAuthorizationURL generates an authorization URL with the given redirect URI, permissions, guildID, disableGuildSelect & scopes. State is automatically generated.
func (c *Client) GenerateAuthorizationURL(redirectURI string, permissions discord.Permissions, guildID snowflake.ID, disableGuildSelect bool, scopes ...discord.OAuth2Scope) string {
	authURL, _ := c.GenerateAuthorizationURLState(redirectURI, permissions, guildID, disableGuildSelect, scopes...)
	return authURL
}

// GenerateAuthorizationURLState generates an authorization URL with the given redirect URI, permissions, guildID, disableGuildSelect & scopes. State is automatically generated & returned.
func (c *Client) GenerateAuthorizationURLState(redirectURI string, permissions discord.Permissions, guildID snowflake.ID, disableGuildSelect bool, scopes ...discord.OAuth2Scope) (string, string) {
	state := c.StateController.NewState(redirectURI)
	values := discord.QueryValues{
		"client_id":     c.ID,
		"redirect_uri":  redirectURI,
		"response_type": "code",
		"scope":         discord.JoinScopes(scopes),
		"state":         state,
	}

	if permissions != discord.PermissionsNone {
		values["permissions"] = permissions
	}
	if guildID != 0 {
		values["guild_id"] = guildID
	}
	if disableGuildSelect {
		values["disable_guild_select"] = true
	}
	return discord.AuthorizeURL(values), state
}

// StartSession starts a new Session with the given authorization code & state.
func (c *Client) StartSession(code string, state string, opts ...rest.RequestOpt) (Session, *discord.IncomingWebhook, error) {
	redirectURI := c.StateController.UseState(state)
	if redirectURI == "" {
		return Session{}, nil, ErrStateNotFound
	}
	accessToken, err := c.Rest.GetAccessToken(c.ID, c.Secret, code, redirectURI, opts...)
	if err != nil {
		return Session{}, nil, err
	}

	return sessionFromTokenResponse(*accessToken), accessToken.Webhook, nil
}

// RefreshSession refreshes the given Session with the refresh token.
func (c *Client) RefreshSession(session Session, opts ...rest.RequestOpt) (Session, error) {
	accessToken, err := c.Rest.RefreshAccessToken(c.ID, c.Secret, session.RefreshToken, opts...)
	if err != nil {
		return Session{}, err
	}
	return sessionFromTokenResponse(*accessToken), nil
}

// VerifySession checks if the Session is expired and refreshes it if needed.
// It returns the new Session, a bool if the Session was expired and an error if something went wrong.
func (c *Client) VerifySession(session Session, opts ...rest.RequestOpt) (Session, bool, error) {
	if session.Expired() {
		newSession, err := c.RefreshSession(session, opts...)
		return newSession, true, err
	}
	return session, false, nil
}

// GetUser returns the discord.OAuth2User associated with the given Session. Fields filled in the struct depend on the Session.Scopes.
func (c *Client) GetUser(session Session, opts ...rest.RequestOpt) (*discord.OAuth2User, error) {
	if err := checkSession(session, discord.OAuth2ScopeIdentify); err != nil {
		return nil, err
	}
	return c.Rest.GetCurrentUser(session.AccessToken, opts...)
}

// GetMember returns the discord.Member associated with the given Session in a specific guild.
func (c *Client) GetMember(session Session, guildID snowflake.ID, opts ...rest.RequestOpt) (*discord.Member, error) {
	if err := checkSession(session, discord.OAuth2ScopeGuildsMembersRead); err != nil {
		return nil, err
	}
	return c.Rest.GetCurrentMember(session.AccessToken, guildID, opts...)
}

// GetGuilds returns the discord.OAuth2Guild(s) the user is a member of. This requires the discord.OAuth2ScopeGuilds scope in the Session.
func (c *Client) GetGuilds(session Session, opts ...rest.RequestOpt) ([]discord.OAuth2Guild, error) {
	if err := checkSession(session, discord.OAuth2ScopeGuilds); err != nil {
		return nil, err
	}
	return c.Rest.GetCurrentUserGuilds(session.AccessToken, 0, 0, 0, false, opts...)
}

// GetConnections returns the discord.Connection(s) the user has connected. This requires the discord.OAuth2ScopeConnections scope in the Session.
func (c *Client) GetConnections(session Session, opts ...rest.RequestOpt) ([]discord.Connection, error) {
	if err := checkSession(session, discord.OAuth2ScopeConnections); err != nil {
		return nil, err
	}
	return c.Rest.GetCurrentUserConnections(session.AccessToken, opts...)
}

// GetApplicationRoleConnection returns the discord.ApplicationRoleConnection for the given application. This requires the discord.OAuth2ScopeRoleConnectionsWrite scope in the Session.
func (c *Client) GetApplicationRoleConnection(session Session, applicationID snowflake.ID, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error) {
	if err := checkSession(session, discord.OAuth2ScopeRoleConnectionsWrite); err != nil {
		return nil, err
	}
	return c.Rest.GetCurrentUserApplicationRoleConnection(session.AccessToken, applicationID, opts...)
}

// UpdateApplicationRoleConnection updates the discord.ApplicationRoleConnection for the given application. This requires the discord.OAuth2ScopeRoleConnectionsWrite scope in the Session.
func (c *Client) UpdateApplicationRoleConnection(session Session, applicationID snowflake.ID, update discord.ApplicationRoleConnectionUpdate, opts ...rest.RequestOpt) (*discord.ApplicationRoleConnection, error) {
	if err := checkSession(session, discord.OAuth2ScopeRoleConnectionsWrite); err != nil {
		return nil, err
	}
	return c.Rest.UpdateCurrentUserApplicationRoleConnection(session.AccessToken, applicationID, update, opts...)
}

func checkSession(session Session, scope discord.OAuth2Scope) error {
	if session.Expired() {
		return ErrSessionExpired
	}
	if !discord.HasScope(scope, session.Scopes...) {
		return ErrMissingOAuth2Scope(scope)
	}
	return nil
}

func sessionFromTokenResponse(accessToken discord.AccessTokenResponse) Session {
	return Session{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: accessToken.RefreshToken,
		Scopes:       accessToken.Scope,
		TokenType:    accessToken.TokenType,
		Expiration:   time.Now().Add(accessToken.ExpiresIn * time.Second),
	}
}
