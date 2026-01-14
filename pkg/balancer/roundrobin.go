package balancer

import (
	"reflect"
	"sync/atomic"
)

type RoundRobin struct {
	// current is an ever-increasing counter.
	// We use atomic operations so RoundRobin is safe for concurrent calls.
	current uint64
}

// NextPeer returns the next node in round-robin order.
//
// Accepts:
//   - slice: []T
//   - pointer to slice: *[]T
//
// Returns (error, peer). If error != nil, peer is nil.
func (r *RoundRobin) NextPeer(nodes interface{}) (error, interface{}) {
	if nodes == nil {
		return ErrInputNotSlice, nil
	}

	v := reflect.ValueOf(nodes)

	// Allow pointer to slice: *[]T
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ErrInputNotSlice, nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return ErrInputNotSlice, nil
	}

	total := v.Len()
	if total == 0 {
		return ErrNoAvaliableNode, nil
	}

	next := r.NextIndex(total)
	return nil, v.Index(next).Interface()
}

// NextIndex returns the next index in [0, total).
// It is safe for concurrent usage.
func (r *RoundRobin) NextIndex(total int) int {
	if total <= 0 {
		return 0
	}

	// Use fetch-and-increment pattern, starting at index 0.
	// atomic.AddUint64 returns the incremented value, so subtract 1 to get previous.
	n := atomic.AddUint64(&r.current, 1) - 1
	return int(n % uint64(total))
}
