package store

import (
	"github.com/appleboy/gin-jwt/v3/core"
)

// StoreType represents the type of store to create
type StoreType string

const (
	// MemoryStore represents an in-memory token store
	MemoryStore StoreType = "memory"
	// RedisStore represents a Redis-based token store
	RedisStore StoreType = "redis"
)

// Config holds the configuration for creating a token store
type Config struct {
	Type  StoreType    // Type of store to create (memory or redis)
	Redis *RedisConfig // Redis configuration (only used when Type is RedisStore)
}

// DefaultConfig returns a default configuration with memory store
func DefaultConfig() *Config { _ = "STUB: not implemented"; return nil }

// NewMemoryConfig creates a configuration for memory store
func NewMemoryConfig() *Config { _ = "STUB: not implemented"; return nil }

// NewRedisConfig creates a configuration for Redis store
func NewRedisConfig(redisConfig *RedisConfig) *Config { _ = "STUB: not implemented"; return nil }

// Factory provides methods to create different types of token stores
type Factory struct{}

// NewFactory creates a new store factory
func NewFactory() *Factory {
	_ = "STUB: not implemented"

	// CreateStore creates a token store based on the provided configuration
	return nil
}

func (f *Factory) CreateStore(config *Config) (core.TokenStore, error) {
	_ = "STUB: not implemented"
	return *new(core.TokenStore), nil
}

// Convenience functions for creating stores

// NewStore creates a token store with the given configuration
func NewStore(config *Config) (core.TokenStore, error) {
	_ = "STUB: not implemented"
	return *new(core.TokenStore), nil
}

// NewMemoryStore creates a new in-memory token store
func NewMemoryStore() core.TokenStore { _ = "STUB: not implemented"; return *new(core.TokenStore) }

// NewRedisStore creates a new Redis token store with the given configuration
func NewRedisStore(config *RedisConfig) (core.TokenStore, error) {
	_ = "STUB: not implemented"
	return *new(core.TokenStore), nil
}

// MustNewStore creates a token store with the given configuration and panics on error
func MustNewStore(config *Config) core.TokenStore {
	_ = "STUB: not implemented"
	return *new(core.TokenStore)
}

// MustNewMemoryStore creates a new in-memory token store (never fails)
func MustNewMemoryStore() core.TokenStore {
	_ = "STUB: not implemented"
	return *

	// MustNewRedisStore creates a new Redis token store and panics on error
	new(core.TokenStore)
}

func MustNewRedisStore(config *RedisConfig) core.TokenStore {
	_ = "STUB: not implemented"
	return *new(core.TokenStore)
}
