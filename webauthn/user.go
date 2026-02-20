package webauthn

import (
	"encoding/base64"
	"fmt"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// SimpleUser implements the WebAuthnUser interface
type SimpleUser struct {
	id          string
	username    string
	displayName string
	credentials []webauthn.Credential
}

// NewSimpleUser creates a new simple user
func NewSimpleUser(id, username, displayName string) *SimpleUser {
	return &SimpleUser{
		id:          id,
		username:    username,
		displayName: displayName,
		credentials: make([]webauthn.Credential, 0),
	}
}

// WebAuthnID returns the user's ID as bytes
func (u *SimpleUser) WebAuthnID() []byte {
	return []byte(u.id)
}

// WebAuthnName returns the user's username
func (u *SimpleUser) WebAuthnName() string {
	return u.username
}

// WebAuthnDisplayName returns the user's display name
func (u *SimpleUser) WebAuthnDisplayName() string {
	return u.displayName
}

// WebAuthnIcon returns the user's icon URL (empty for simple implementation)
func (u *SimpleUser) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials returns the user's credentials
func (u *SimpleUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

// GetUsername returns the username
func (u *SimpleUser) GetUsername() string {
	return u.username
}

// GetDisplayName returns the display name
func (u *SimpleUser) GetDisplayName() string {
	return u.displayName
}

// SetCredential adds or updates a credential
func (u *SimpleUser) SetCredential(cred webauthn.Credential) {
	// Check if credential already exists and update it
	for i, existingCred := range u.credentials {
		if string(existingCred.ID) == string(cred.ID) {
			u.credentials[i] = cred
			return
		}
	}
	// If not found, add as new credential
	u.credentials = append(u.credentials, cred)
}

// GetCredentials returns all credentials
func (u *SimpleUser) GetCredentials() []webauthn.Credential {
	return u.credentials
}

// CredentialExcludeList returns credential descriptor list for registration
func (u *SimpleUser) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         "public-key",
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}
	return credentialExcludeList
}

// SimpleUserStore implements UserStoreInterface using in-memory storage
type SimpleUserStore struct {
	users map[string]*SimpleUser
}

// NewSimpleUserStore creates a new simple user store
func NewSimpleUserStore() *SimpleUserStore {
	return &SimpleUserStore{
		users: make(map[string]*SimpleUser),
	}
}

// GetUser retrieves a user by ID
func (s *SimpleUserStore) GetUser(userID string) (WebAuthnUser, error) {
	user, exists := s.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	return user, nil
}

// CreateUser creates a new user
func (s *SimpleUserStore) CreateUser(userID, username, displayName string) (WebAuthnUser, error) {
	user := NewSimpleUser(userID, username, displayName)
	s.users[userID] = user
	return user, nil
}

// UpdateUser updates an existing user
func (s *SimpleUserStore) UpdateUser(user WebAuthnUser) error {
	simpleUser, ok := user.(*SimpleUser)
	if !ok {
		return fmt.Errorf("invalid user type")
	}
	s.users[simpleUser.id] = simpleUser
	return nil
}

// SimpleSessionStore implements SessionStoreInterface using in-memory storage
type SimpleSessionStore struct {
	sessions map[string]*webauthn.SessionData
}

// NewSimpleSessionStore creates a new simple session store
func NewSimpleSessionStore() *SimpleSessionStore {
	return &SimpleSessionStore{
		sessions: make(map[string]*webauthn.SessionData),
	}
}

// StoreSession stores session data
func (s *SimpleSessionStore) StoreSession(sessionID string, data *webauthn.SessionData) error {
	s.sessions[sessionID] = data
	return nil
}

// GetSession retrieves session data
func (s *SimpleSessionStore) GetSession(sessionID string) (*webauthn.SessionData, error) {
	data, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}
	return data, nil
}

// DeleteSession removes session data
func (s *SimpleSessionStore) DeleteSession(sessionID string) error {
	delete(s.sessions, sessionID)
	return nil
}

// Utility functions

// CredentialToBase64 converts credential ID to base64
func CredentialToBase64(credID []byte) string {
	return base64.URLEncoding.EncodeToString(credID)
}

// CredentialFromBase64 converts base64 string back to credential ID
func CredentialFromBase64(credStr string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(credStr)
}