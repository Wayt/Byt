package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestScheme(t *testing.T) {

	requests := []struct {
		URL            string
		Header         http.Header
		ExpectedScheme string
	}{
		{
			URL:            "https://localhost:8080",
			Header:         nil,
			ExpectedScheme: "https",
		},

		{
			URL:            "http://localhost:8080",
			Header:         nil,
			ExpectedScheme: "http",
		},
		{
			URL: "http://localhost",
			Header: http.Header{
				"X-Forwarded-Proto": []string{"https"},
			},
			ExpectedScheme: "https",
		},
	}

	for _, req := range requests {
		r := httptest.NewRequest("GET", req.URL, nil)
		if req.Header != nil {
			r.Header = req.Header
		}

		scheme := RequestScheme(r)
		if scheme != req.ExpectedScheme {
			t.Errorf("%+v: invalid scheme, got %s", req, scheme)
		}
	}
}

func TestDetectContentType(t *testing.T) {

	files := []struct {
		Filename            string
		Data                []byte
		ExpectedContentType string
	}{
		{
			Filename: "all.min.css",
			Data: []byte(`.center {
    text-align: center;
    color: red;
}`),
			ExpectedContentType: "text/css",
		},
		{
			Filename: "script.js",
			Data: []byte(`function(){
			console.log('hello world');
		}();`),
			ExpectedContentType: "text/javascript",
		},
		{
			Filename:            "image.png",
			Data:                []byte(`dummy`),
			ExpectedContentType: "image/png",
		},
		{
			Filename:            "dummy",
			Data:                []byte(`hello world`),
			ExpectedContentType: "text/plain; charset=utf-8",
		},
	}

	for _, f := range files {
		contentType := DetectContentType(f.Filename, f.Data)

		if contentType != f.ExpectedContentType {
			t.Errorf("%s: invalid content type, expected %s, got %s", f.Filename, f.ExpectedContentType, contentType)
		}
	}
}
