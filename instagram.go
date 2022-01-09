package ig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	profileParams      = "fields=id,username,account_type"
	mediaParams        = "fields=id,username,caption,media_type,media_url,permalink,thumbnail_url,timestamp"
	exchangeTokenParam = "grant_type=ig_exchange_token"
)

// Client instagram connection representation.
type Client struct {
	accessToken string
	client      *http.Client
	baseURL     string
}

func (c *Client) WithTransport(newClient *http.Client) {
	c.client = newClient
}

// NewClient creates a new client with a given accessToken and clientSecret.
func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		baseURL:     "https://graph.instagram.com",
		client:      &http.Client{},
	}
}

// Entry represents an instagram media post.
type Entry struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Caption      string `json:"caption"`
	MediaType    string `json:"media_type"`
	MediaURL     string `json:"media_url"`
	Permalink    string `json:"permalink"`
	ThumbnailURL string `json:"thumbnail_url"`
	Timestamp    string `json:"timestamp"`
}

func (e *Entry) String() string {
	return fmt.Sprintf(
		"%T{ID: %s,Username: %s,Caption: %s,MediaType: %s,MediaURL: %s,Permalink: %s,ThumbnailURL: %s,Timestamp:%s}",
		e,
		e.ID,
		e.Username,
		e.Caption,
		e.MediaType,
		e.MediaURL,
		e.Permalink,
		e.ThumbnailURL,
		e.Timestamp,
	)
}

// Tags returns all identified hashtags in caption.
func (e Entry) Tags() []string {
	var result []string

	arr := strings.Split(e.Caption, " ")
	for _, s := range arr {
		if strings.HasPrefix(s, "#") {
			for _, h := range strings.Split(s, "#") {
				if h == "" {
					continue
				}
				result = append(result, h)
			}
		}
	}

	return result
}

// Profile represents an instagram profile, it's retrieved with GetProfile.
type Profile struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	AccountType string `json:"account_type"`
	MediaCount  string `json:"media_count"`
}

func (p *Profile) String() string {
	return fmt.Sprintf(
		"%T{ID: %s, Username: %s, AccountType: %s, MediaCount: %s}",
		p,
		p.ID,
		p.Username,
		p.AccountType,
		p.MediaCount,
	)
}

// LongLiveToken represents a response for long live token request.
type LongLiveToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (t *LongLiveToken) String() string {
	return fmt.Sprintf(
		"%T{AccessToken: %s, TokenType: %s, ExpiresIn: %d}",
		t,
		t.AccessToken,
		t.TokenType,
		t.ExpiresIn,
	)
}

type mediaResp struct {
	Data []Entry `json:"data"`
}

// GetMedia fetches media from a user configured within a Client, returns an array of Entries,
// error if something goes wrong with communication.
func (c *Client) GetMedia() ([]Entry, error) {
	var (
		url        = buildURL(c.baseURL, "/me/media", c.accessToken, mediaParams)
		bytes, err = c.fetch(url)
	)

	if err != nil {
		return nil, err
	}

	var resp mediaResp
	if err = json.Unmarshal(bytes, &resp); err != nil {
		return nil, errors.Wrapf(err, "unable to Unmarshal json response: %s", string(bytes))
	}

	return resp.Data, nil
}

// GetProfile fetches instagram user profile information.
func (c *Client) GetProfile() (*Profile, error) {
	var (
		url        = buildURL(c.baseURL, "/me", c.accessToken, profileParams)
		bytes, err = c.fetch(url)
	)

	if err != nil {
		return nil, err
	}

	var p Profile
	if err = json.Unmarshal(bytes, &p); err != nil {
		return nil, errors.Wrapf(err, "unable to Unmarshal json response: %s", string(bytes))
	}

	return &p, nil
}

// LongLivedToken returns long live toke from instagram graph, secret corresponds to Instagram app secret.
func (c *Client) LongLivedToken(secret string) (*LongLiveToken, error) {
	var (
		secretParam = fmt.Sprintf("client_secret=%s", secret)
		url         = buildURL(c.baseURL, "/access_token", c.accessToken, exchangeTokenParam, secretParam)
		bytes, err  = c.fetch(url)

		token LongLiveToken
	)

	err = json.Unmarshal(bytes, &token)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal long-live-token")
	}

	return &token, nil
}

func (c *Client) fetch(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "cannot fetch profile")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	if err = resp.Body.Close(); err != nil {
		return nil, errors.Wrap(err, "unable to close body")
	}

	return body, nil
}

func buildURL(base, path, token string, extraParams ...string) string {
	var params string

	if len(extraParams) > 0 {
		params = fmt.Sprintf("&%s", strings.Join(extraParams, "&"))
	}

	return fmt.Sprintf("%s%s?access_token=%s%s", base, path, token, params)
}
