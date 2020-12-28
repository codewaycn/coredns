// +build fuzz

package cache

import (
	"coredns/plugin/pkg/fuzz"
)

// Fuzz fuzzes cache.
func Fuzz(data []byte) int {
	return fuzz.Do(New(), data)
}
