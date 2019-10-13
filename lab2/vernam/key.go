package vernam

import (
	"math/rand"
	"time"
)

func CreateKey(n uint64) []byte {
	rand.Seed(time.Now().UnixNano())
	key := make([]byte, n)
	for i := uint64(0); i < n; i++ {
		key[i] = byte(rand.Intn(255))
	}
	return key
}
