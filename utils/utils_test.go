package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	assert := assert.New(t)
	assert.True(Exists(RootDir() + "/main.go"))
	assert.False(Exists(RootDir() + "/main.og"))
}

// BenchmarkExists-12    	  243614	      4351 ns/op	     552 B/op	       5 allocs/op
func BenchmarkExists(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Exists(RootDir() + "/main.go")
	}
}
