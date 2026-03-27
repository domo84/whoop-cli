package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// callbackServer listens on localhost for the OAuth2 redirect.
type callbackServer struct {
	port   int
	codeCh chan string
	errCh  chan error
	state  string
	server *http.Server
}

func newCallbackServer(port int, state string) *callbackServer {
	s := &callbackServer{
		port:  port,
		state: state,
		codeCh: make(chan string, 1),
		errCh:  make(chan error, 1),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/callback/whoop", s.handleCallback)
	s.server = &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s
}

// Start begins listening and returns the redirect_uri to embed in the auth URL.
func (s *callbackServer) Start() (string, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.port))
	if err != nil {
		return "", fmt.Errorf("starting callback server on port %d: %w", s.port, err)
	}
	go s.server.Serve(ln) //nolint:errcheck
	return fmt.Sprintf("http://localhost:%d/callback/whoop", s.port), nil
}

// Wait blocks until the authorization code is received or ctx is cancelled.
func (s *callbackServer) Wait(ctx context.Context) (string, error) {
	select {
	case code := <-s.codeCh:
		return code, nil
	case err := <-s.errCh:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// Shutdown stops the HTTP server gracefully.
func (s *callbackServer) Shutdown(ctx context.Context) {
	shutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	s.server.Shutdown(shutCtx) //nolint:errcheck
}

func (s *callbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Validate state to prevent CSRF.
	if got := q.Get("state"); got != s.state {
		s.errCh <- fmt.Errorf("state mismatch: expected %q got %q", s.state, got)
		http.Error(w, "state mismatch", http.StatusBadRequest)
		return
	}

	if errMsg := q.Get("error"); errMsg != "" {
		desc := q.Get("error_description")
		s.errCh <- fmt.Errorf("oauth2 error: %s — %s", errMsg, desc)
		http.Error(w, "authentication failed", http.StatusBadRequest)
		return
	}

	code := q.Get("code")
	if code == "" {
		s.errCh <- fmt.Errorf("no code in callback")
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	s.codeCh <- code
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, successHTML)
}

const successHTML = `<!DOCTYPE html>
<html>
<head><title>Whoop CLI Authentication</title>
<style>body{font-family:sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0;background:#f0f4f8;}
.box{text-align:center;padding:2rem;border-radius:8px;background:white;box-shadow:0 2px 8px rgba(0,0,0,.1);}
h1{color:#2d9c52;}p{color:#555;}</style></head>
<body><div class="box"><h1>Authentication successful</h1>
<p>You may close this tab and return to your terminal.</p></div></body>
</html>`
