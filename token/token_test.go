package token

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type getTokenTransport struct {
	StatusCode int
	Body       string
}

func (t *getTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: t.StatusCode,
		Body:       ioutil.NopCloser(strings.NewReader(t.Body)),
	}
	return resp, nil
}

func TestToken_Get(t *testing.T) {
	httpClient := &http.Client{}
	httpClient.Transport = &getTokenTransport{
		StatusCode: http.StatusOK,
		Body:       `{"access_token": "test_token"}`,
	}
	cli := NewClient("http://localhost", "id", "secret", WithHTTPClient(httpClient))

	accessToken, err := cli.GetAccessToken(context.Background(), "code")
	if err != nil {
		t.Fatal(err)
	}

	if accessToken != "test_token" {
		t.Errorf("Expect token test_token, got %v", accessToken)
	}
}
