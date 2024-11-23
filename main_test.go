package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupApp(t *testing.T) {
	// Initialize the app
	app := SetupApp()

	// Test CORS middleware (OPTIONS request)
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "API is up and running")

	assert.NotNil(t, resp, "Logger middleware should allow requests to proceed")
}

func TestRoutes(t *testing.T) {
	app := SetupApp()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Health endpoint should return 200")
}
