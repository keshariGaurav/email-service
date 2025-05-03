package retry

import (
	"log"
	"time"
)

// RetryWithBackoff retries a function with exponential backoff
func RetryWithBackoff(fn func() error, maxAttempts int) error {
	var lastErr error
	backoff := time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := fn(); err != nil {
			lastErr = err
			log.Printf("Attempt %d/%d failed: %v", attempt, maxAttempts, err)
			if attempt == maxAttempts {
				break
			}
			time.Sleep(backoff)
			backoff *= 2
			continue
		}
		return nil
	}
	return lastErr
}
