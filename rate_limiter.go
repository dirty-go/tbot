package tbot

import (
	"math"
	"sync"
	"time"
)

// rateLimiter is a simple token-bucket limiter. A nil receiver is a
// no-op so callers can leave the limiter unset to disable throttling.
type rateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	burst      float64
	rate       float64
	lastRefill time.Time

	now   func() time.Time
	sleep func(time.Duration)
}

// newRateLimiter returns a limiter that admits up to rps requests per
// second with a maximum burst of burst tokens. rps <= 0 or burst <= 0
// disables throttling and returns nil.
func newRateLimiter(rps float64, burst int) *rateLimiter {
	if rps <= 0 || burst <= 0 {
		return nil
	}
	return &rateLimiter{
		tokens:     float64(burst),
		burst:      float64(burst),
		rate:       rps,
		lastRefill: time.Now(),
		now:        time.Now,
		sleep:      time.Sleep,
	}
}

// wait blocks until a token is available, then consumes it.
func (rl *rateLimiter) wait() {
	if rl == nil {
		return
	}
	for {
		rl.mu.Lock()
		now := rl.now()
		if elapsed := now.Sub(rl.lastRefill).Seconds(); elapsed > 0 {
			rl.tokens = math.Min(rl.burst, rl.tokens+elapsed*rl.rate)
			rl.lastRefill = now
		}
		if rl.tokens >= 1 {
			rl.tokens--
			rl.mu.Unlock()
			return
		}
		deficit := 1 - rl.tokens
		wait := time.Duration(deficit / rl.rate * float64(time.Second))
		rl.mu.Unlock()
		if wait <= 0 {
			wait = time.Millisecond
		}
		rl.sleep(wait)
	}
}
