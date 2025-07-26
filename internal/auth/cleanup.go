package auth

import (
	"context"
	"log"
	"time"
)

func (a *AuthService) StartTokenCleanup(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour) // Cleanup daily
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := a.CleanupExpiredTokens(); err != nil {
				log.Printf("Error cleaning up expired tokens: %v", err)
			} else {
				log.Println("Successfully cleaned up expired tokens")
			}
		}
	}
}
