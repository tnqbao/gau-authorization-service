package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

const (
	RefreshTokenBlacklistPrefix = "refresh_token_blacklist:"
	RefreshTokenTTLPrefix       = "refresh_token_ttl:"
)

// Generate a new UUID for refresh token
func (r *Repository) GenerateRefreshTokenID(_ context.Context) (uuid.UUID, error) {
	return uuid.New(), nil
}

// Blacklist a refresh token UUID
func (r *Repository) BlacklistRefreshTokenID(ctx context.Context, id uuid.UUID) error {
	key := RefreshTokenBlacklistPrefix + id.String()
	if _, err := r.CacheDB.Set(ctx, key, 1, 0).Result(); err != nil {
		return err
	}
	return nil
}

// Check if a refresh token UUID is blacklisted
func (r *Repository) IsRefreshTokenBlacklisted(ctx context.Context, id uuid.UUID) (bool, error) {
	key := RefreshTokenBlacklistPrefix + id.String()
	exists, err := r.CacheDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// Blacklist refresh token with TTL
func (r *Repository) BlacklistRefreshTokenIDWithTTL(ctx context.Context, id uuid.UUID, ttl time.Duration) error {
	pipe := r.CacheDB.TxPipeline()

	// Set the ID in the blacklist
	blacklistKey := RefreshTokenBlacklistPrefix + id.String()
	pipe.Set(ctx, blacklistKey, 1, ttl)

	// Set TTL tracking key
	ttlKey := RefreshTokenTTLPrefix + id.String()
	pipe.Set(ctx, ttlKey, 1, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

// Clean up expired blacklist entries
func (r *Repository) CleanupBlacklistEntries(ctx context.Context) error {
	// Scan for all blacklist entries
	var cursor uint64
	cleanedCount := 0

	for {
		keys, nextCursor, err := r.CacheDB.Scan(ctx, cursor, RefreshTokenBlacklistPrefix+"*", 1000).Result()
		if err != nil {
			return fmt.Errorf("scan error: %w", err)
		}

		for _, key := range keys {
			// Check if the key still exists (hasn't expired)
			exists, err := r.CacheDB.Exists(ctx, key).Result()
			if err != nil {
				log.Printf("Error checking key existence %s: %v\n", key, err)
				continue
			}

			// If key doesn't exist, it was auto-expired by Redis TTL
			if exists == 0 {
				cleanedCount++
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	log.Printf("Blacklist cleanup complete. %d entries were auto-expired by Redis TTL.\n", cleanedCount)
	return nil
}

// Remove from blacklist (for manual cleanup)
func (r *Repository) RemoveFromBlacklist(ctx context.Context, id uuid.UUID) error {
	pipe := r.CacheDB.TxPipeline()

	blacklistKey := RefreshTokenBlacklistPrefix + id.String()
	ttlKey := RefreshTokenTTLPrefix + id.String()

	pipe.Del(ctx, blacklistKey)
	pipe.Del(ctx, ttlKey)

	_, err := pipe.Exec(ctx)
	return err
}

// Legacy method compatibility - no longer needed but kept for interface compatibility
func (r *Repository) GetBit(_ context.Context, _ string, _ int64) (int64, error) {
	return 0, fmt.Errorf("GetBit is deprecated for UUID-based tokens")
}
