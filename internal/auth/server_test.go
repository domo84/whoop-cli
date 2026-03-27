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

	srv := newCallbackServer(port, state)
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
	srv := newCallbackServer(port, "expected-state")
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
	srv := newCallbackServer(port, "state")
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
	srv := newCallbackServer(port, "state")
	_, err := srv.Start()
	require.NoError(t, err)
	defer srv.Shutdown(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err = srv.Wait(ctx)
	require.Error(t, err)
}
