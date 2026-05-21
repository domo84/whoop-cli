package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"

	"github.com/domo84/whoop-cli/internal/config"
)

const (
	authURL  = "https://api.prod.whoop.com/oauth/oauth2/auth"
	tokenURL = "https://api.prod.whoop.com/oauth/oauth2/token"
	revokeURL = "https://api.prod.whoop.com/v2/user/oauth/revokeUserOauthAccess"
)

var defaultScopes = []string{
	"read:profile",
	"read:body_measurement",
	"read:cycles",
	"read:sleep",
	"read:recovery",
	"read:workout",
	"offline",
}

// OAuthFlow orchestrates the full authorization code + PKCE flow.
type OAuthFlow interface {
	Login(ctx context.Context) (*oauth2.Token, error)
	Revoke(ctx context.Context, token *oauth2.Token) error
}

type pkceFlow struct {
	cfg        *config.Config
	oauthConf  *oauth2.Config
	httpClient *http.Client
}

// NewPKCEFlow creates an OAuthFlow using the given Config for client credentials.
func NewPKCEFlow(cfg *config.Config) OAuthFlow {
	return NewPKCEFlowWithHTTPClient(cfg, http.DefaultClient)
}

// NewPKCEFlowWithHTTPClient creates an OAuthFlow with a custom HTTP client (for testing).
func NewPKCEFlowWithHTTPClient(cfg *config.Config, httpClient *http.Client) OAuthFlow {
	return &pkceFlow{
		cfg: cfg,
		oauthConf: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Scopes:       defaultScopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  authURL,
				TokenURL: tokenURL,
			},
		},
		httpClient: httpClient,
	}
}

func (f *pkceFlow) Login(ctx context.Context) (*oauth2.Token, error) {
	pkce, err := NewPKCEChallenge()
	if err != nil {
		return nil, fmt.Errorf("generating PKCE challenge: %w", err)
	}

	state, err := generateState()
	if err != nil {
		return nil, fmt.Errorf("generating state: %w", err)
	}

	srv, err := newCallbackServer(f.cfg.RedirectURI, state)
	if err != nil {
		return nil, err
	}
	redirectURI, err := srv.Start()
	if err != nil {
		return nil, err
	}
	defer srv.Shutdown(ctx)

	// Build authorization URL with PKCE params.
	authURLStr := f.oauthConf.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", pkce.Challenge),
		oauth2.SetAuthURLParam("code_challenge_method", pkce.Method),
		oauth2.SetAuthURLParam("redirect_uri", redirectURI),
	)

	fmt.Printf("Callback URL: %s\n", redirectURI)
	fmt.Printf("Make sure this exact URI is registered as a redirect in your WHOOP developer app.\n\n")
	fmt.Printf("Opening browser for authentication...\n")
	fmt.Printf("If your browser does not open automatically, visit:\n\n  %s\n\n", authURLStr)
	openBrowser(authURLStr)

	code, err := srv.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("waiting for authorization: %w", err)
	}

	// Inject custom HTTP client into context for token exchange.
	ctx = context.WithValue(ctx, oauth2.HTTPClient, f.httpClient)

	token, err := f.oauthConf.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", pkce.Verifier),
		oauth2.SetAuthURLParam("redirect_uri", redirectURI),
	)
	if err != nil {
		return nil, fmt.Errorf("exchanging code for token: %w", err)
	}
	return token, nil
}

func (f *pkceFlow) Revoke(ctx context.Context, token *oauth2.Token) error {
	body, _ := json.Marshal(map[string]string{"token": token.AccessToken})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, revokeURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building revoke request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("revoking token: %w", err)
	}
	defer resp.Body.Close()
	// A 200 or 204 means success; ignore other non-critical errors.
	return nil
}

// NewTokenSource returns an oauth2.TokenSource that automatically refreshes
// the access token using the refresh token and persists the new token to disk.
func NewTokenSource(ctx context.Context, cfg *config.Config, token *oauth2.Token, store TokenStore) oauth2.TokenSource {
	oauthConf := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
	base := oauthConf.TokenSource(ctx, token)
	return &persistingTokenSource{base: base, store: store}
}

// persistingTokenSource wraps an oauth2.TokenSource and saves refreshed tokens to disk.
type persistingTokenSource struct {
	base  oauth2.TokenSource
	store TokenStore
}

func (p *persistingTokenSource) Token() (*oauth2.Token, error) {
	token, err := p.base.Token()
	if err != nil {
		return nil, err
	}
	// Save whenever the token may have been refreshed.
	_ = p.store.Save(token)
	return token, nil
}

// openBrowser attempts to open the given URL in the default browser.
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start() //nolint:errcheck
}
