package linenotify

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type AuthorizationClient struct {
	ClientID     string
	RedirectURI  string
	ResponseMode string
	State        string
}

type AuthorizationResponse struct {
	Code             string
	State            string
	Error            string
	ErrorDescription string
}

func NewAuthorization(clientID, redirectURI string) (*AuthorizationClient, error) {
	state, err := generateHash()
	if err != nil {
		return nil, err
	}

	return &AuthorizationClient{
		ClientID:     clientID,
		RedirectURI:  redirectURI,
		ResponseMode: "form_post",
		State:        state,
	}, nil
}

func (c *AuthorizationClient) RequestURL() (string, error) {
	u, err := url.Parse("https://notify-bot.line.me/oauth/authorize")
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Add("response_type", "code")
	v.Add("client_id", c.ClientID)
	v.Add("redirect_uri", c.RedirectURI)
	v.Add("scope", "notify")
	v.Add("state", c.State)
	v.Add("response_mode", c.ResponseMode)
	u.RawQuery = v.Encode()

	return u.String(), nil
}

func (c *AuthorizationClient) Redirect(w http.ResponseWriter, req *http.Request) error {
	urlStr, err := c.RequestURL()
	if err != nil {
		return err
	}
	http.Redirect(w, req, urlStr, http.StatusFound)
	return nil
}

func generateHash() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func ParseAuthorization(r *http.Request) (*AuthorizationResponse, error) {
	resp := &AuthorizationResponse{
		Code:             r.FormValue("code"),
		State:            r.FormValue("state"),
		Error:            r.FormValue("error"),
		ErrorDescription: r.FormValue("error_description"),
	}

	if resp.Error != "" {
		return resp, errors.Errorf("authorize failure. %s", resp.ErrorDescription)
	}
	return resp, nil
}
