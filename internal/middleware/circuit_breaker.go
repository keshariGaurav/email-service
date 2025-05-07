package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
    cb *gobreaker.CircuitBreaker
}

func NewCircuitBreaker() *CircuitBreaker {
    settings := gobreaker.Settings{
        Name:        "email-service-breaker",
        MaxRequests: 3,                     // Number of requests allowed in half-open state
        Interval:    60 * time.Second,      // Time window for counting failures
        Timeout:     30 * time.Second,      // How long to wait before attempting recovery
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            // Only trip if:
            // 1. We have at least 5 requests AND
            // 2. The failure ratio is >= 60%
            if counts.Requests < 5 {
                return false
            }
            failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
            return failureRatio >= 0.6
        },
        OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
            // Log state changes for debugging
            println("Circuit Breaker state changed from:", from.String(), "to:", to.String())
        },
    }
    

    return &CircuitBreaker{
        cb: gobreaker.NewCircuitBreaker(settings),
    }
}

func CircuitBreakerMiddleware() fiber.Handler {
    breaker := NewCircuitBreaker()
    
    return func(c *fiber.Ctx) error {
        resultChan := make(chan error, 1)
        timeoutCtx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
        defer cancel() // Ensure context is cancelled to prevent goroutine leak
        
        // Execute request with timeout
        go func() {
            defer close(resultChan) // Ensure channel is closed
            
            _, err := breaker.cb.Execute(func() (interface{}, error) {
                // Create a channel to track Next() completion
                done := make(chan error, 1)
                
                go func() {
                    defer close(done)
                    done <- c.Next()
                }()
                
                // Wait for Next() or context timeout
                select {
                case err := <-done:
                    if err != nil {
                        return nil, fmt.Errorf("service error: %w", err)
                    }
                    
                    // Check response status after Next() completes
                    status := c.Response().StatusCode()
                    if status >= 500 {
                        return nil, fmt.Errorf("service returned %d status code", status)
                    }
                    
                    // Count authentication failures for internal services
                    if status == 401 {
                        return nil, fmt.Errorf("authentication failure in internal service call")
                    }
                    
                    return nil, nil
                    
                case <-timeoutCtx.Done():
                    return nil, fmt.Errorf("request timed out")
                }
            })
            
            resultChan <- err
        }()

        // Wait for result or timeout
        select {
        case err := <-resultChan:
            if err != nil {
                if err == gobreaker.ErrOpenState {
                    return c.Status(503).JSON(fiber.Map{
                        "error": "Circuit breaker is open",
                        "details": "Service is temporarily unavailable due to multiple failures",
                        "retry_after": "30s",
                    })
                }
                
                if err.Error() == "request timed out" {
                    return c.Status(504).JSON(fiber.Map{
                        "error": "Service timeout",
                        "details": "Request exceeded 5 second timeout",
                    })
                }
                
                return c.Status(500).JSON(fiber.Map{
                    "error": "Service error",
                    "details": err.Error(),
                })
            }
            
            return nil
            
        case <-timeoutCtx.Done():
            return c.Status(504).JSON(fiber.Map{
                "error": "Service timeout",
                "details": "Request exceeded 5 second timeout",
            })
        }
    }
}