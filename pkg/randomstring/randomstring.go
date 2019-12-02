package randomstring

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func Generator(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return string(b)
}
