package retry

import (
	"log"
	"time"
)

func RetryWithBackoff(fn func() error, attempts int) {
	backoff := time.Second

	for i := 0; i < attempts; i++ {
		err := fn()
		if err == nil {
			log.Printf("Email sent successfully after %d retries.", i+1)
			return
		}
		log.Printf("Retry %d failed: %v", i+1, err)
		time.Sleep(backoff)
		backoff *= 2
	}
}
