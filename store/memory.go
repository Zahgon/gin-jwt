package store

import (
	"context"
	"sync"
	"time"

	"github.com/appleboy/gin-jwt/v3/core"
)

var _ core.TokenStore = &InMemoryRefreshTokenStore{}

// InMemoryRefreshTokenStore provides a simple in-memory refresh token store
// This implementation is thread-safe and suitable for single-instance applications
// For distributed systems, consider using Redis or database-based implementations
type InMemoryRefreshTokenStore struct {
	tokens map[string]*core.RefreshTokenData
	mu     sync.RWMutex
}

// NewInMemoryRefreshTokenStore creates a new in-memory refresh token store
func NewInMemoryRefreshTokenStore() *InMemoryRefreshTokenStore {
	_ = "STUB: not implemented"
	return nil
}

// Set stores a refresh token with associated user data and expiration
func (s *InMemoryRefreshTokenStore) Set(
	ctx context.Context,
	token string,
	userData any,
	expiry time.Time,
) error {
	_ = "STUB: not implemented"
	return nil
}

// Get retrieves user data associated with a refresh token
func (s *InMemoryRefreshTokenStore) Get(ctx context.Context, token string) (any, error) {
	_ = "STUB: not implemented"
	return *new(any), nil
}

// Clean up expired token

// Delete removes a refresh token from storage
func (s *InMemoryRefreshTokenStore) Delete(ctx context.Context, token string) error {
	_ = "STUB: not implemented"
	return nil

	// No error for empty token deletion
}

// Cleanup removes expired tokens and returns the number of tokens cleaned up
func (s *InMemoryRefreshTokenStore) Cleanup(ctx context.Context) (int, error) {
	_ = "STUB: not implemented"
	return 0, nil
}

// Count returns the total number of active refresh tokens
func (s *InMemoryRefreshTokenStore) Count(ctx context.Context) (int, error) {
	_ = "STUB: not implemented"
	return 0, nil
}

// GetAll returns all active tokens (for debugging/monitoring purposes)
// Note: This method is not part of the RefreshTokenStorer interface
// and should be used carefully in production environments
func (s *InMemoryRefreshTokenStore) GetAll() map[string]*core.RefreshTokenData {
	_ = "STUB: not implemented"
	return nil
}

// Create a copy to prevent external modifications

// Clear removes all tokens from the store (useful for testing)
// Note: This method is not part of the RefreshTokenStorer interface
func (s *InMemoryRefreshTokenStore) Clear() { _ = "STUB: not implemented"; return }
