package balancer

import "strings"

// Balancer defines a load-balancing strategy.
//
// NOTE: Signature is kept for backward compatibility with existing implementations
// (e.g., RoundRobin) in other files.
type Balancer interface {
	// NextPeer selects the next peer from nodes.
	// Returns (error, peer). If error != nil, peer should be ignored.
	NextPeer(nodes interface{}) (error, interface{})
}

const (
	// StrategyRoundRobin is the canonical strategy name (normalized form).
	StrategyRoundRobin = "roundrobin"
)

// normalizeStrategy makes strategy matching tolerant to case, spaces, '-' and '_'.
// Examples that become "roundrobin":
// "RoundRobin", "round_robin", "round-robin", " rr ", etc.
func normalizeStrategy(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "-", "")
	return s
}

// NewBalancer returns a Balancer based on the given strategy.
// Unknown strategies will fall back to RoundRobin for safety.
func NewBalancer(strategy string) Balancer {
	switch normalizeStrategy(strategy) {
	case "", StrategyRoundRobin, "rr":
		return &RoundRobin{}
	default:
		// Fallback behavior: keep service running with a sane default.
		return &RoundRobin{}
	}
}
