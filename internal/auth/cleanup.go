package auth

import (
	"context"
	"time"
)

func (a *AuthService) StartTokenCleanup(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour) // Cleanup daily
	defer ticker.Stop()

	a.logger.Info("Token cleanup service started", "interval", "24h")

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("Token cleanup service stopped")
			return
		case <-ticker.C:
			a.logger.Info("Starting scheduled token cleanup")
			if err := a.CleanupExpiredTokens(ctx); err != nil {
				a.logger.LogError(ctx, err, "Failed to cleanup expired tokens")
			} else {
				a.logger.Info("Token cleanup completed successfully")
			}
		}
	}
}
