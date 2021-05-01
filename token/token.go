package token

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type (
	// Client represents api client that get access token
	Client struct {
		HTTPClient   *http.Client
		RedirectURI  string
		ClientID     string
		ClientSecret string
	}

	response struct {
		AccessToken string `json:"access_token"`
	}

	// Option with client
	Option func(*Client)
)

// NewClient returns Client
func NewClient(redirectURI, clientID, clientSecret string, opts ...Option) *Client {
	c := &Client{
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
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// GetAccessToken returns access token that published by line notify
func (c *Client) GetAccessToken(ctx context.Context, code string) (string, error) {
	v := url.Values{}
	v.Add("grant_type", "authorization_code")
	v.Add("code", code)
	v.Add("redirect_uri", c.RedirectURI)
	v.Add("client_id", c.ClientID)
	v.Add("client_secret", c.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://notify-bot.line.me/oauth/token", strings.NewReader(v.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		respToken := &response{}
		if err := json.NewDecoder(resp.Body).Decode(respToken); err != nil {
			return "", fmt.Errorf("failed to decode toke response: %w", err)
		}
		return respToken.AccessToken, nil
	}
	return "", errors.Errorf("failed to get access token. status:%v", resp.Status)
}
