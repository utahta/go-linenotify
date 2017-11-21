package auth

import (
	"net/http"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	req, err := New("id", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}

	if req.ResponseMode != "form_post" {
		t.Errorf("Expect form_post, got %v", req.ResponseMode)
	}

	if len(req.State) != 36 {
		t.Errorf("Expect state length 44, got %v", len(req.State))
	}
}

func TestClient_RequestURL(t *testing.T) {
	req, err := New("id", "http://localhost/linenotify_test")
	if err != nil {
		t.Fatal(err)
	}

	urlStr, err := req.RequestURL()
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		key         string
		expectValue string
	}{

		{"response_type", "code"},
		{"client_id", "id"},
		{"redirect_uri", "http://localhost/linenotify_test"},
		{"scope", "notify"},
		{"state", req.State},
		{"response_mode", "form_post"},
	}
	for _, test := range tests {
		v := u.Query().Get(test.key)
		if test.expectValue != v {
			t.Errorf("Expect %v:%v, got %v", test.key, test.expectValue, v)
		}
	}
}

func TestParseAuthorize(t *testing.T) {

	tests := []struct {
		urlValues   url.Values
		expectError bool
	}{
		{url.Values{
			"code":              []string{"code_value"},
			"state":             []string{"state_value"},
			"error":             []string{""},
			"error_description": []string{""},
		}, false},

		{url.Values{
			"code":              []string{"code_value"},
			"state":             []string{"state_value"},
			"error":             []string{"error_value"},
			"error_description": []string{"error_description_value"},
		}, true},
	}

	for _, test := range tests {
		req := &http.Request{Form: test.urlValues}
		resp, err := ParseRequest(req)

		if test.expectError {
			if err == nil {
				t.Errorf("Expect error, got nil. urlValues:%#v", test.urlValues)
			}

			if resp.ErrorDescription != "error_description_value" {
				t.Errorf("Expect 'error_description_value', got %v", resp.ErrorDescription)
			}
		} else {
			if err != nil {
				t.Errorf("Expect nil, got error. urlValues:%#v", test.urlValues)
			}

			if resp.ErrorDescription != "" {
				t.Errorf("Expect '', got %v", resp.ErrorDescription)
			}
		}
	}
}
