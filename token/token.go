package token

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type (
	// Client interface
	Client interface {
		GetAccessToken(string) (string, error)
	}

	client struct {
		HTTPClient   *http.Client
		RedirectURI  string
		ClientID     string
		ClientSecret string
	}

	response struct {
		AccessToken string `json:"access_token"`
	}

	// Option with client
	Option func(*client)
)

// New returns token get client
func New(redirectURI, clientID, clientSecret string, opts ...Option) Client {
	c := &client{
		HTTPClient:   http.DefaultClient,
		RedirectURI:  redirectURI,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	for _, o := range opts {
		o(c)
	}
	return c
}

// WithHTTPClient set the http client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *client) {
		c.HTTPClient = httpClient
	}
}

// GetAccessToken returns access token that published by line notify
func (c *client) GetAccessToken(code string) (string, error) {
	v := url.Values{}
	v.Add("grant_type", "authorization_code")
	v.Add("code", code)
	v.Add("redirect_uri", c.RedirectURI)
	v.Add("client_id", c.ClientID)
	v.Add("client_secret", c.ClientSecret)

	resp, err := c.HTTPClient.PostForm("https://notify-bot.line.me/oauth/token", v)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		respToken := &response{}
		if err := json.NewDecoder(resp.Body).Decode(respToken); err != nil {
			return "", err
		}
		return respToken.AccessToken, nil
	}
	return "", errors.Errorf("failed to get access token. status:%v", resp.Status)
}
