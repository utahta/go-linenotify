package linenotify

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type TokenClient struct {
	HTTPClient   *http.Client
	Code         string
	RedirectURI  string
	ClientID     string
	ClientSecret string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

// NewToken returns *TokenClient
func NewToken(code, redirectURI, clientID, clientSecret string) *TokenClient {
	return &TokenClient{
		HTTPClient:   http.DefaultClient,
		Code:         code,
		RedirectURI:  redirectURI,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// Get returns access token
func (c *TokenClient) Get() (string, error) {
	v := url.Values{}
	v.Add("grant_type", "authorization_code")
	v.Add("code", c.Code)
	v.Add("redirect_uri", c.RedirectURI)
	v.Add("client_id", c.ClientID)
	v.Add("client_secret", c.ClientSecret)

	resp, err := c.HTTPClient.PostForm("https://notify-bot.line.me/oauth/token", v)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		respToken := &TokenResponse{}
		if err := json.NewDecoder(resp.Body).Decode(respToken); err != nil {
			return "", err
		}
		return respToken.AccessToken, nil
	}
	return "", errors.Errorf("failed to get access token. status:%v", resp.Status)
}
