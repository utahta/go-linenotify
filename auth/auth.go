package auth

import (
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	// Client represents LINE Notify authorization
	Client struct {
		ClientID     string
		RedirectURI  string
		ResponseMode string
		State        string
	}

	// AuthorizeResponse represents LINE Notify authorize response
	AuthorizeResponse struct {
		Code             string
		State            string
		Error            string
		ErrorDescription string
	}
)

// New returns Client
func New(clientID, redirectURI string) (*Client, error) {
	randomID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Client{
		ClientID:     clientID,
		RedirectURI:  redirectURI,
		ResponseMode: "form_post",
		State:        randomID.String(),
	}, nil
}

// RequestURL builds request url with parameters
func (c *Client) RequestURL() (string, error) {
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

// Redirect redirect to request url
func (c *Client) Redirect(w http.ResponseWriter, req *http.Request) error {
	urlStr, err := c.RequestURL()
	if err != nil {
		return err
	}
	http.Redirect(w, req, urlStr, http.StatusFound)
	return nil
}

// ParseRequest parses authorize request
func ParseRequest(r *http.Request) (*AuthorizeResponse, error) {
	resp := &AuthorizeResponse{
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
