package webauthn

import (
	"testing"
)

func TestNewSimpleUserStore(t *testing.T) {
	store := NewSimpleUserStore()
	if store == nil {
		t.Fatal("NewSimpleUserStore() returned nil")
	}
	
	if len(store.users) != 0 {
		t.Errorf("Expected empty user store, got %d users", len(store.users))
	}
}

func TestNewSimpleSessionStore(t *testing.T) {
	store := NewSimpleSessionStore()
	if store == nil {
		t.Fatal("NewSimpleSessionStore() returned nil")
	}
	
	if len(store.sessions) != 0 {
		t.Errorf("Expected empty session store, got %d sessions", len(store.sessions))
	}
}

func TestSimpleUser(t *testing.T) {
	user := NewSimpleUser("test-id", "test-user", "Test User")
	
	if user.WebAuthnID() == nil {
		t.Error("WebAuthnID() returned nil")
	}
	
	if string(user.WebAuthnID()) != "test-id" {
		t.Errorf("Expected WebAuthnID 'test-id', got '%s'", string(user.WebAuthnID()))
	}
	
	if user.WebAuthnName() != "test-user" {
		t.Errorf("Expected WebAuthnName 'test-user', got '%s'", user.WebAuthnName())
	}
	
	if user.WebAuthnDisplayName() != "Test User" {
		t.Errorf("Expected WebAuthnDisplayName 'Test User', got '%s'", user.WebAuthnDisplayName())
	}
	
	if user.WebAuthnIcon() != "" {
		t.Errorf("Expected empty WebAuthnIcon, got '%s'", user.WebAuthnIcon())
	}
	
	if len(user.WebAuthnCredentials()) != 0 {
		t.Errorf("Expected 0 credentials, got %d", len(user.WebAuthnCredentials()))
	}
}

func TestSimpleUserStore(t *testing.T) {
	store := NewSimpleUserStore()
	
	// Test CreateUser
	user, err := store.CreateUser("test-id", "test-user", "Test User")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	
	// Test GetUser
	retrievedUser, err := store.GetUser("test-id")
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	
	if retrievedUser.GetUsername() != "test-user" {
		t.Errorf("Expected username 'test-user', got '%s'", retrievedUser.GetUsername())
	}
	
	// Test UpdateUser
	err = store.UpdateUser(user)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}
	
	// Test GetUser for non-existent user
	_, err = store.GetUser("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
}

func TestWebAuthnConfig(t *testing.T) {
	userStore := NewSimpleUserStore()
	sessionStore := NewSimpleSessionStore()
	
	// Test missing required fields
	_, err := New(&WebAuthnConfig{})
	if err == nil {
		t.Error("Expected error for empty config, got nil")
	}
	
	// Test with minimum required fields
	client, err := New(&WebAuthnConfig{
		ClientID:      "test-client",
		ClientSecret:  "test-secret",
		RedirectURI:   "http://localhost:8080/callback",
		RPDisplayName: "Test App",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:8080"},
		UserStore:     userStore,
		SessionStore:  sessionStore,
	})
	
	if err != nil {
		t.Fatalf("New() failed with valid config: %v", err)
	}
	
	if client == nil {
		t.Fatal("New() returned nil client")
	}
}

func TestCredentialToBase64(t *testing.T) {
	testData := []byte("test-credential-id")
	encoded := CredentialToBase64(testData)
	
	decoded, err := CredentialFromBase64(encoded)
	if err != nil {
		t.Fatalf("CredentialFromBase64 failed: %v", err)
	}
	
	if string(decoded) != string(testData) {
		t.Errorf("Expected '%s', got '%s'", string(testData), string(decoded))
	}
}

func TestCredentialFromBase64(t *testing.T) {
	// Test invalid base64
	_, err := CredentialFromBase64("invalid-base64!@#")
	if err == nil {
		t.Error("Expected error for invalid base64, got nil")
	}
}