package util

import (
	"math/rand"
	"time"
)

var (
	r = rand.New(rand.NewSource(time.Now().Unix()))
)

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func RandCount(max int64) int64 {
	rand.Seed(int64(time.Now().Nanosecond()))
	if max == 0 {
		return 0
	}
	b := rand.Int63n(max)
	return b
}
