package tbot

import (
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiter_DisabledForNonPositive(t *testing.T) {
	cases := []struct {
		name  string
		rps   float64
		burst int
	}{
		{"zero rps", 0, 10},
		{"negative rps", -1, 10},
		{"zero burst", 10, 0},
		{"negative burst", 10, -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if rl := newRateLimiter(tc.rps, tc.burst); rl != nil {
				t.Fatalf("expected nil limiter for rps=%v burst=%d, got %+v", tc.rps, tc.burst, rl)
			}
		})
	}
}

func TestRateLimiter_NilWaitIsNoop(t *testing.T) {
	var rl *rateLimiter
	rl.wait()
}

func TestRateLimiter_ConsumesBurstWithoutSleeping(t *testing.T) {
	rl := newRateLimiter(1, 5)
	rl.sleep = func(time.Duration) {
		t.Fatal("sleep should not be called while burst tokens remain")
	}
	for range 5 {
		rl.wait()
	}
}

func TestRateLimiter_BlocksWhenEmptyAndSleepsExpectedDuration(t *testing.T) {
	rl := newRateLimiter(2, 1) // 1 token initially, refill 2/sec → ~500ms per token
	rl.wait()                  // consumes the initial token

	// freeze "now" so the sleep duration is deterministic.
	frozen := time.Unix(0, 0)
	rl.now = func() time.Time { return frozen }

	var slept time.Duration
	called := 0
	rl.sleep = func(d time.Duration) {
		called++
		slept += d
		// pretend the sleep advances time enough to fully refill one token.
		frozen = frozen.Add(d)
		rl.lastRefill = frozen.Add(-d) // keep refill math consistent
	}

	rl.wait()

	if called == 0 {
		t.Fatal("expected sleep to be called when bucket is empty")
	}
	// 1 token at 2 tokens/sec → ~500ms. Allow a healthy margin.
	if slept < 400*time.Millisecond || slept > 2*time.Second {
		t.Fatalf("unexpected total sleep %v", slept)
	}
}

func TestRateLimiter_ConcurrentSafety(t *testing.T) {
	rl := newRateLimiter(1000, 1000) // generous so wait never blocks meaningfully
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 20 {
				rl.wait()
			}
		}()
	}
	wg.Wait()
}
