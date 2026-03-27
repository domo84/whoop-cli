package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// APIServer is a test HTTP server that serves canned JSON responses.
type APIServer struct {
	*httptest.Server
	mux      *http.ServeMux
	mu       sync.Mutex
	callLog  map[string]int // method+path -> count
}

// NewAPIServer creates a new test server. Remember to call srv.Close() in cleanup.
func NewAPIServer(t *testing.T) *APIServer {
	t.Helper()
	mux := http.NewServeMux()
	srv := &APIServer{
		mux:     mux,
		callLog: make(map[string]int),
	}
	srv.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + " " + r.URL.Path
		srv.mu.Lock()
		srv.callLog[key]++
		srv.mu.Unlock()
		mux.ServeHTTP(w, r)
	}))
	return srv
}

// Handle registers a handler for method+path that returns the given status and JSON body.
func (s *APIServer) Handle(method, path string, status int, body any) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			json.NewEncoder(w).Encode(body) //nolint:errcheck
		}
	})
}

// HandleFunc registers a raw handler for more complex test scenarios.
func (s *APIServer) HandleFunc(method, path string, fn http.HandlerFunc) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fn(w, r)
	})
}

// HandleAny registers a handler for any HTTP method on a path.
func (s *APIServer) HandleAny(path string, fn http.HandlerFunc) {
	s.mux.HandleFunc(path, fn)
}

// CallCount returns how many times method+path was requested.
func (s *APIServer) CallCount(method, path string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.callLog[method+" "+path]
}

// AssertCalled asserts that method+path was called exactly n times.
func (s *APIServer) AssertCalled(t *testing.T, method, path string, times int) {
	t.Helper()
	got := s.CallCount(method, path)
	if got != times {
		t.Errorf("expected %s %s to be called %d time(s), got %d", method, path, times, got)
	}
}
