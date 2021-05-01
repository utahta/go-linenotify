package linenotify

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/sync/errgroup"
)

type notifyRoundTripper struct {
	resp *http.Response
	err  error
}

func (t *notifyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t.resp.Request = req
	return t.resp, t.err
}

func TestClient_Notify(t *testing.T) {
	c := NewClient()
	statusOK := `{"status":200,"message":"ok"}`
	statusUnauthorized := `{"status":401,"message":"invalid access token"}`
	tests := []struct {
		resp           *http.Response
		imageThumbnail string
		imageFullsize  string
		image          io.Reader
		expectedErr    error
		explain        string
	}{
		{&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(statusOK))}, "", "", nil, nil, "ok: message"},
		{&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(statusOK))}, "image.jpg", "image.jpg", nil, nil, "ok: message image url"},
		{&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(statusOK))}, "", "", bytes.NewReader([]byte("image file")), nil, "ok: message image"},
		{&http.Response{StatusCode: http.StatusUnauthorized, Body: ioutil.NopCloser(strings.NewReader(statusUnauthorized))}, "", "", nil, ErrNotifyInvalidAccessToken, "ng: message"},
	}

	for _, test := range tests {
		c.HTTPClient.Transport = &notifyRoundTripper{resp: test.resp}

		_, err := c.Notify(context.Background(), "token", "test", test.imageThumbnail, test.imageFullsize, test.image)
		if err != test.expectedErr {
			t.Errorf("%v err:%v", test.explain, err)
		}
	}
}

func TestClient_requestBodyWithImage(t *testing.T) {
	c := NewClient()

	c.HTTPClient.Transport = &notifyRoundTripper{
		resp: &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(""))},
		err:  nil,
	}

	body, contentType, err := c.requestBodyWithImage("test", bytes.NewReader([]byte("image file")))
	if err != nil {
		t.Fatal(err)
	}
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(buf), "image file") {
		t.Errorf("Expected buffer image file, got %s", string(buf))
	}

	if !strings.Contains(contentType, "multipart/form-data;") {
		t.Errorf("Expected contentType, got %s", contentType)
	}

	// for data race
	eg := &errgroup.Group{}
	c.HTTPClient.Transport = &notifyRoundTripper{
		resp: &http.Response{Body: ioutil.NopCloser(strings.NewReader("image file"))},
		err:  nil,
	}

	for i := 0; i < 30; i++ {
		eg.Go(func() error {
			_, _, err := c.requestBodyWithImage("test", bytes.NewReader([]byte("image binary")))
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
