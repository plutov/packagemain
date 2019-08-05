package main

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter .
type IPRateLimiter struct {
	visitors map[string]*visitor
	mu       *sync.RWMutex
	// rate
	r           rate.Limit
	burstTokens int
}

// NewIPRateLimiter .
func NewIPRateLimiter(r rate.Limit, burstTokens int) *IPRateLimiter {
	i := &IPRateLimiter{
		visitors:    make(map[string]*visitor),
		mu:          &sync.RWMutex{},
		r:           r,
		burstTokens: burstTokens,
	}

	go i.CleanupVisitors()

	return i
}

// AddVisitor creates a new rate limiter and adds it to the visitors map, using the
// IP address as the key.
func (i *IPRateLimiter) AddVisitor(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.burstTokens)

	// Include the current time when creating a new visitor.
	i.visitors[ip] = &visitor{limiter, time.Now()}

	return limiter
}

// GetVisitor retrieves and returns the rate limiter for the current visitor if it
// already exists. Otherwise call the addVisitor function to add a
// new entry to the map.
func (i *IPRateLimiter) GetVisitor(ip string) *rate.Limiter {
	i.mu.Lock()
	v, exists := i.visitors[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddVisitor(ip)
	}

	i.mu.Unlock()

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

// CleanupVisitors checks the map for visitors that haven't been seen for
// more than 3 minutes and deletes the entries.
func (i *IPRateLimiter) CleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		i.mu.Lock()
		for ip, v := range i.visitors {
			if time.Now().Sub(v.lastSeen) > 3*time.Minute {
				delete(i.visitors, ip)
			}
		}
		i.mu.Unlock()
	}
}
