package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port
}

func TestCallbackServer_ValidCallback(t *testing.T) {
	port := freePort(t)
	state := "test-state-value"

	srv, err := newCallbackServer(fmt.Sprintf("http://localhost:%d/callback/whoop", port), state)
	require.NoError(t, err)
	redirectURI, err := srv.Start()
	require.NoError(t, err)
	defer srv.Shutdown(context.Background())

	go func() {
		time.Sleep(20 * time.Millisecond)
		callbackURL := fmt.Sprintf("%s?code=mycode&state=%s", redirectURI, url.QueryEscape(state))
		resp, err := http.Get(callbackURL) //nolint:noctx
		if err == nil {
			resp.Body.Close()
		}
	}()

	code, err := srv.Wait(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "mycode", code)
}

func TestCallbackServer_StateMismatch(t *testing.T) {
	port := freePort(t)
	srv, err := newCallbackServer(fmt.Sprintf("http://localhost:%d/callback/whoop", port), "expected-state")
	require.NoError(t, err)
	redirectURI, err := srv.Start()
	require.NoError(t, err)
	defer srv.Shutdown(context.Background())

	go func() {
		time.Sleep(20 * time.Millisecond)
		callbackURL := fmt.Sprintf("%s?code=mycode&state=wrong-state", redirectURI)
		resp, err := http.Get(callbackURL) //nolint:noctx
		if err == nil {
			resp.Body.Close()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = srv.Wait(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "state mismatch")
}

func TestCallbackServer_OAuthError(t *testing.T) {
	port := freePort(t)
	srv, err := newCallbackServer(fmt.Sprintf("http://localhost:%d/callback/whoop", port), "state")
	require.NoError(t, err)
	redirectURI, err := srv.Start()
	require.NoError(t, err)
	defer srv.Shutdown(context.Background())

	go func() {
		time.Sleep(20 * time.Millisecond)
		callbackURL := fmt.Sprintf("%s?error=access_denied&error_description=User+denied&state=state", redirectURI)
		resp, err := http.Get(callbackURL) //nolint:noctx
		if err == nil {
			resp.Body.Close()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = srv.Wait(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "access_denied")
}

func TestCallbackServer_ContextCancelled(t *testing.T) {
	port := freePort(t)
	srv, err := newCallbackServer(fmt.Sprintf("http://localhost:%d/callback/whoop", port), "state")
	require.NoError(t, err)
	_, err = srv.Start()
	require.NoError(t, err)
	defer srv.Shutdown(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err = srv.Wait(ctx)
	require.Error(t, err)
}

func TestCallbackServer_InvalidURI(t *testing.T) {
	cases := []struct {
		name string
		uri  string
		want string
	}{
		{"missing port", "http://localhost/callback", "explicit port is required"},
		{"missing path", "http://localhost:8282", "path is required"},
		{"wrong scheme", "ftp://localhost:8282/callback", "scheme must be http or https"},
		{"empty", "", "scheme must be http or https"},
		{"no host", "http:///callback", "host is required"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newCallbackServer(tc.uri, "state")
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}
