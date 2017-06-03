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
}

var (
	ErrNotifyInvalidAccessToken = errors.New("Invalid access token.")
)

// https://notify-bot.line.me/doc/
func New() *Client {
	return &Client{HTTPClient: http.DefaultClient}
}

func (c *Client) Notify(token, message, imageThumbnail, imageFullsize string, image io.Reader) error {
	if image != nil {
		return c.NotifyWithImage(token, message, image)
	}
	return c.NotifyWithImageURL(token, message, imageThumbnail, imageFullsize)
}

func (c *Client) NotifyWithImage(token, message string, image io.Reader) error {
	body, contentType, err := c.requestBodyWithImage(message, image)
	if err != nil {
		return err
	}
	return c.notify(token, message, body, contentType)
}

func (c *Client) NotifyWithImageURL(token, message, imageThumbnail, imageFullsize string) error {
	body, contentType, err := c.requestBody(message, imageThumbnail, imageFullsize)
	if err != nil {
		return err
	}
	return c.notify(token, message, body, contentType)
}

func (c *Client) notify(token, message string, body io.Reader, contentType string) error {
	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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
