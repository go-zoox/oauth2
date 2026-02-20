package webauthn

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/go-zoox/oauth2"
)

// WebAuthnConfig holds the configuration for WebAuthn authentication
type WebAuthnConfig struct {
	// Basic OAuth2-like configuration
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`

	// WebAuthn specific configuration
	RPDisplayName string `json:"rp_display_name"` // Display name for your app
	RPID          string `json:"rp_id"`           // Domain for your app
	RPOrigins     []string `json:"rp_origins"`    // Allowed origins

	// User store interface
	UserStore UserStoreInterface `json:"-"`

	// Session store interface
	SessionStore SessionStoreInterface `json:"-"`

	// Optional: Timeout for authentication (in milliseconds)
	Timeout uint64 `json:"timeout"`
}

// UserStoreInterface defines the interface for user storage operations
type UserStoreInterface interface {
	GetUser(userID string) (WebAuthnUser, error)
	CreateUser(userID, username, displayName string) (WebAuthnUser, error)
	UpdateUser(user WebAuthnUser) error
}

// SessionStoreInterface defines the interface for session storage operations
type SessionStoreInterface interface {
	StoreSession(sessionID string, data *webauthn.SessionData) error
	GetSession(sessionID string) (*webauthn.SessionData, error)
	DeleteSession(sessionID string) error
}

// WebAuthnUser implements the webauthn.User interface
type WebAuthnUser interface {
	webauthn.User
	GetUsername() string
	GetDisplayName() string
	SetCredential(cred webauthn.Credential)
	GetCredentials() []webauthn.Credential
}

// WebAuthnClientInterface provides additional WebAuthn-specific methods
type WebAuthnClientInterface interface {
	oauth2.Client
	BeginRegistration(userID, username, displayName string) (*protocol.PublicKeyCredentialCreationOptions, string, error)
	FinishRegistration(userID, sessionID string, credentialResponse []byte) error
	BeginLogin(userID string) (*protocol.PublicKeyCredentialRequestOptions, string, error)
	FinishLogin(userID, sessionID string, assertionResponse []byte) error
}

// client implements the OAuth2 client interface for WebAuthn
type client struct {
	config   *WebAuthnConfig
	webauthn *webauthn.WebAuthn
}

// New creates a new WebAuthn OAuth2 client
func New(cfg *WebAuthnConfig) (WebAuthnClientInterface, error) {
	if cfg.RPDisplayName == "" {
		return nil, fmt.Errorf("webauthn: rp_display_name is required")
	}
	if cfg.RPID == "" {
		return nil, fmt.Errorf("webauthn: rp_id is required")
	}
	if len(cfg.RPOrigins) == 0 {
		return nil, fmt.Errorf("webauthn: rp_origins is required")
	}
	if cfg.UserStore == nil {
		return nil, fmt.Errorf("webauthn: user_store is required")
	}
	if cfg.SessionStore == nil {
		return nil, fmt.Errorf("webauthn: session_store is required")
	}

	// Set default timeout
	if cfg.Timeout == 0 {
		cfg.Timeout = 60000 // 60 seconds
	}

	// Create WebAuthn configuration
	wconfig := &webauthn.Config{
		RPDisplayName: cfg.RPDisplayName,
		RPID:          cfg.RPID,
		RPOrigins:     cfg.RPOrigins,
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce: true,
				Timeout: time.Duration(cfg.Timeout) * time.Millisecond,
			},
			Registration: webauthn.TimeoutConfig{
				Enforce: true,
				Timeout: time.Duration(cfg.Timeout) * time.Millisecond,
			},
		},
	}

	// Create WebAuthn instance
	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		return nil, fmt.Errorf("webauthn: failed to create webauthn instance: %v", err)
	}

	return &client{
		config:   cfg,
		webauthn: webAuthn,
	}, nil
}

// Authorize starts the WebAuthn authentication process
// For WebAuthn, this can start either registration or login flow
// The state parameter is used to determine the flow and store session data
func (c *client) Authorize(state string, callback func(loginUrl string)) {
	// Parse state to determine if this is registration or login
	params, _ := url.ParseQuery(state)
	flow := params.Get("flow")
	userID := params.Get("user_id")

	if flow == "" {
		flow = "login" // default to login
	}

	// Create the authorization URL with WebAuthn-specific parameters
	authURL := fmt.Sprintf("%s?flow=%s&user_id=%s&state=%s", 
		c.config.RedirectURI, flow, userID, state)

	callback(authURL)
}

// Callback handles the WebAuthn response
// For WebAuthn, this will process either registration or login completion
func (c *client) Callback(code, state string, cb func(user *oauth2.User, token *oauth2.Token, err error)) {
	// Parse the request to get the WebAuthn response
	// Note: In a real implementation, you would get this from the HTTP request
	// This is a simplified version for demonstration
	
	params, err := url.ParseQuery(state)
	if err != nil {
		cb(nil, nil, fmt.Errorf("webauthn: invalid state parameter: %v", err))
		return
	}

	flow := params.Get("flow")
	userID := params.Get("user_id")

	switch flow {
	case "registration":
		c.handleRegistrationCallback(userID, code, state, cb)
	case "login":
		c.handleLoginCallback(userID, code, state, cb)
	default:
		cb(nil, nil, fmt.Errorf("webauthn: unknown flow: %s", flow))
	}
}

// handleRegistrationCallback processes WebAuthn registration completion
func (c *client) handleRegistrationCallback(userID, code, state string, cb func(user *oauth2.User, token *oauth2.Token, err error)) {
	// Get the stored session data
	_, err := c.config.SessionStore.GetSession(state)
	if err != nil {
		cb(nil, nil, fmt.Errorf("webauthn: failed to get session: %v", err))
		return
	}

	// Get the user
	user, err := c.config.UserStore.GetUser(userID)
	if err != nil {
		cb(nil, nil, fmt.Errorf("webauthn: failed to get user: %v", err))
		return
	}

	// Parse the credential response (this would come from the HTTP request in practice)
	// For now, we'll simulate a successful registration
	
	// Create a dummy credential for demonstration
	// In practice, you would call webauthn.FinishRegistration here
	
	// Create OAuth2 user representation
	oauthUser := &oauth2.User{
		ID:       userID,
		Username: user.GetUsername(),
		Email:    user.GetUsername(), // Assuming username is email
		Nickname: user.GetDisplayName(),
	}

	// Create OAuth2 token representation
	token := &oauth2.Token{
		AccessToken:  generateAccessToken(),
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
	}

	// Clean up session
	c.config.SessionStore.DeleteSession(state)

	cb(oauthUser, token, nil)
}

// handleLoginCallback processes WebAuthn login completion
func (c *client) handleLoginCallback(userID, code, state string, cb func(user *oauth2.User, token *oauth2.Token, err error)) {
	// Get the stored session data
	_, err := c.config.SessionStore.GetSession(state)
	if err != nil {
		cb(nil, nil, fmt.Errorf("webauthn: failed to get session: %v", err))
		return
	}

	// Get the user
	user, err := c.config.UserStore.GetUser(userID)
	if err != nil {
		cb(nil, nil, fmt.Errorf("webauthn: failed to get user: %v", err))
		return
	}

	// Parse the assertion response (this would come from the HTTP request in practice)
	// For now, we'll simulate a successful login
	
	// In practice, you would call webauthn.FinishLogin here
	
	// Create OAuth2 user representation
	oauthUser := &oauth2.User{
		ID:       userID,
		Username: user.GetUsername(),
		Email:    user.GetUsername(), // Assuming username is email
		Nickname: user.GetDisplayName(),
	}

	// Create OAuth2 token representation
	token := &oauth2.Token{
		AccessToken:  generateAccessToken(),
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
	}

	// Clean up session
	c.config.SessionStore.DeleteSession(state)

	cb(oauthUser, token, nil)
}

// Logout handles WebAuthn logout (mainly cleanup)
func (c *client) Logout(callback func(logoutUrl string)) {
	// WebAuthn doesn't have a traditional logout URL
	// Just redirect to a logout confirmation page
	logoutURL := "/logout"
	callback(logoutURL)
}

// Register handles WebAuthn user registration
func (c *client) Register(callback func(registerUrl string)) {
	// Create registration URL
	registerURL := fmt.Sprintf("%s?flow=registration", c.config.RedirectURI)
	callback(registerURL)
}

// RefreshToken is not applicable to WebAuthn
func (c *client) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return nil, fmt.Errorf("webauthn: refresh token not supported")
}

// Helper functions

func generateAccessToken() string {
	// In practice, you would generate a proper JWT or session token
	return "webauthn_access_token_" + fmt.Sprintf("%d", time.Now().Unix())
}

// WebAuthn-specific helper methods

// BeginRegistration starts the WebAuthn registration ceremony
func (c *client) BeginRegistration(userID, username, displayName string) (*protocol.PublicKeyCredentialCreationOptions, string, error) {
	// Get or create user
	user, err := c.config.UserStore.GetUser(userID)
	if err != nil {
		// Create new user if not exists
		user, err = c.config.UserStore.CreateUser(userID, username, displayName)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create user: %v", err)
		}
	}

	// Begin registration
	options, sessionData, err := c.webauthn.BeginRegistration(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to begin registration: %v", err)
	}

	// Store session
	sessionID := generateSessionID()
	err = c.config.SessionStore.StoreSession(sessionID, sessionData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to store session: %v", err)
	}

	return &options.Response, sessionID, nil
}

// FinishRegistration completes the WebAuthn registration ceremony
func (c *client) FinishRegistration(userID, sessionID string, credentialResponse []byte) error {
	// Get session
	_, err := c.config.SessionStore.GetSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %v", err)
	}

	// Get user
	user, err := c.config.UserStore.GetUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	// In a real implementation, you would parse the HTTP request here
	// For now, we'll simulate a successful registration completion
	// 
	// The credential response would be parsed from the HTTP request like this:
	// parsedResponse, err := protocol.ParseCredentialCreationResponseBody(r.Body)
	// if err != nil {
	//     return fmt.Errorf("failed to parse credential response: %v", err)
	// }
	// 
	// credential, err := c.webauthn.FinishRegistration(user, *session, parsedResponse)
	// if err != nil {
	//     return fmt.Errorf("failed to finish registration: %v", err)
	// }

	// For demo purposes, we'll create a dummy credential
	credential := &webauthn.Credential{
		ID:              []byte("dummy-credential-id"),
		PublicKey:       []byte("dummy-public-key"),
		AttestationType: "none",
		Transport:       []protocol.AuthenticatorTransport{"internal"},
		Flags: webauthn.CredentialFlags{
			UserPresent:    true,
			UserVerified:   true,
			BackupEligible: false,
			BackupState:    false,
		},
		Authenticator: webauthn.Authenticator{
			AAGUID:    []byte("dummy-aaguid"),
			SignCount: 0,
		},
	}

	// Store credential
	user.SetCredential(*credential)
	err = c.config.UserStore.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	// Clean up session
	c.config.SessionStore.DeleteSession(sessionID)

	return nil
}

// BeginLogin starts the WebAuthn login ceremony
func (c *client) BeginLogin(userID string) (*protocol.PublicKeyCredentialRequestOptions, string, error) {
	// Get user
	user, err := c.config.UserStore.GetUser(userID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %v", err)
	}

	// Begin login
	options, sessionData, err := c.webauthn.BeginLogin(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to begin login: %v", err)
	}

	// Store session
	sessionID := generateSessionID()
	err = c.config.SessionStore.StoreSession(sessionID, sessionData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to store session: %v", err)
	}

	return &options.Response, sessionID, nil
}

// FinishLogin completes the WebAuthn login ceremony
func (c *client) FinishLogin(userID, sessionID string, assertionResponse []byte) error {
	// Get session
	_, err := c.config.SessionStore.GetSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %v", err)
	}

	// Get user
	user, err := c.config.UserStore.GetUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	// In a real implementation, you would parse the HTTP request here
	// For now, we'll simulate a successful login completion
	// 
	// The assertion response would be parsed from the HTTP request like this:
	// parsedResponse, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	// if err != nil {
	//     return fmt.Errorf("failed to parse assertion response: %v", err)
	// }
	// 
	// credential, err := c.webauthn.FinishLogin(user, *session, parsedResponse)
	// if err != nil {
	//     return fmt.Errorf("failed to finish login: %v", err)
	// }

	// For demo purposes, we'll update an existing credential
	credentials := user.GetCredentials()
	if len(credentials) == 0 {
		return fmt.Errorf("user has no registered credentials")
	}
	
	credential := credentials[0] // Use first credential
	// Update sign count
	credential.Authenticator.SignCount++

	// Update credential (sign count, etc.)
	user.SetCredential(credential)
	err = c.config.UserStore.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	// Clean up session
	c.config.SessionStore.DeleteSession(sessionID)

	return nil
}

func generateSessionID() string {
	return fmt.Sprintf("webauthn_session_%d", time.Now().UnixNano())
}