package linenotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Client struct {
	HTTPClient *http.Client
	token      string
}

type ClientOption func(*Client)

var (
	ErrNotifyInvalidAccessToken = errors.New("Invalid access token.")
)

// https://notify-bot.line.me/doc/ja/
func New(options ...ClientOption) *Client {
	c := &Client{HTTPClient: http.DefaultClient}

	for _, opt := range options {
		opt(c)
	}
	return c
}

func WithToken(token string) ClientOption {
	return func(c *Client) {
		c.token = token
	}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) Notify(message, imageThumbnail, imageFullsize string, image io.Reader) error {
	if image != nil {
		return c.NotifyWithImage(message, image)
	}
	return c.NotifyWithImageURL(message, imageThumbnail, imageFullsize)
}

func (c *Client) NotifyWithImage(message string, image io.Reader) error {
	body, contentType, err := c.requestBodyWithImage(message, image)
	if err != nil {
		return err
	}
	return c.notify(message, body, contentType)
}

func (c *Client) NotifyWithImageURL(message, imageThumbnail, imageFullsize string) error {
	body, contentType, err := c.requestBody(message, imageThumbnail, imageFullsize)
	if err != nil {
		return err
	}
	return c.notify(message, body, contentType)
}

func (c *Client) notify(message string, body io.Reader, contentType string) error {
	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrNotifyInvalidAccessToken
	}

	if resp.StatusCode != http.StatusOK {
		var data interface{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return err
		}
		root := data.(map[string]interface{})
		return errors.New(root["message"].(string))
	}
	return nil
}

func (c *Client) requestBody(message, imageThumbnail, imageFullsize string) (io.Reader, string, error) {
	v := url.Values{}
	v.Add("message", message)
	if imageThumbnail != "" {
		v.Add("imageThumbnail", imageThumbnail)
	}
	if imageFullsize != "" {
		v.Add("imageFullsize", imageFullsize)
	}
	return strings.NewReader(v.Encode()), "application/x-www-form-urlencoded", nil
}

func (c *Client) requestBodyWithImage(message string, image io.Reader) (io.Reader, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	if err := w.WriteField("message", message); err != nil {
		return nil, "", err
	}

	filename, err := generateHash()
	if err != nil {
		return nil, "", err
	}

	fw, err := w.CreateFormFile("imageFile", filename)
	if err != nil {
		return nil, "", err
	}

	if _, err := io.Copy(fw, image); err != nil {
		return nil, "", err
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return &b, w.FormDataContentType(), nil
}
