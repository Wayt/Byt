package main

import (
	"net/http"
	"path/filepath"
)

// RequestScheme returns the Scheme of an input request.
// It will priorize header value other actual request scheme.
// Defaults to 'http'.
func RequestScheme(r *http.Request) string {

	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}

	if scheme := r.URL.Scheme; scheme != "" {
		return scheme
	}

	return "http"
}

// DetectContentType returns best fiting file content type.
// It priorize file extension other http's `DetectContentType` function.
func DetectContentType(filename string, data []byte) string {

	switch filepath.Ext(filename) {
	case ".js":
		return "text/javascript"
	case ".css":
		return "text/css"
	case ".png":
		return "image/png"
	default:
		return http.DetectContentType(data)
	}
}
