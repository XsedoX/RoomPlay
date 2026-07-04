package cache_cleanup_worker

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func StartCacheCleanupWorker(ctx context.Context, db *sql.DB, interval time.Duration, ttl time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := db.Exec(
				"DELETE FROM cache WHERE created_at_utc < $1",
				time.Now().Add(-ttl).UTC(),
			)
			if err != nil {
				log.Printf("Error cleaning up cache: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
