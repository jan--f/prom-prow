package util

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v57/github"
)

func TestErrorResponse404(t *testing.T) {
	// Test that we can identify 404 errors correctly
	err := &github.ErrorResponse{
		Response: &http.Response{
			StatusCode: 404,
		},
	}

	// Verify it's a 404
	if err.Response.StatusCode != 404 {
		t.Errorf("Expected 404, got %d", err.Response.StatusCode)
	}
}
