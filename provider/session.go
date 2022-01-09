package provider

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/markbates/goth"
	"github.com/pkg/errors"

	"github.com/dvarelap/ig.git"
)

// Session stores data during the auth process with Instagram
type Session struct {
	AuthURL     string    `json:"auth_url"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	Secret      string    `json:"secret"`
}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the Instagram provider.
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}

	return s.AuthURL, nil
}

// Authorize the session with Instagram and return the access token to be stored for future use.
func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)
	token, err := p.config.Exchange(goth.ContextForClient(p.Client()), params.Get("code"))
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("Invalid token received from provider")
	}

	client := ig.NewClient(token.AccessToken)

	ltToken, err := client.LongLivedToken(s.Secret)

	if err != nil {
		return "", err
	}

	s.AccessToken = ltToken.AccessToken
	s.ExpiresAt = time.Now().Add(time.Duration(ltToken.ExpiresIn) * time.Second)

	return token.AccessToken, err
}

// Marshal the session into a string
func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s Session) String() string {
	return s.Marshal()
}

// UnmarshalSession will unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	sess := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(sess)
	return sess, err
}
