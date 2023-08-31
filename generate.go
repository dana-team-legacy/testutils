package testutils

import (
	"math/rand"
	"time"
)

const (
	e2ePrefix        = "e2e"
	randStringLength = 8
)

// GenerateE2EName generates a name for e2e testing
func GenerateE2EName(nm, testPrefix, randPrefix string) string {
	prefix := e2ePrefix + "-" + testPrefix + "-" + randPrefix + "-"
	name := prefix + nm

	return name
}

// RandStr generates a random string
func RandStr() string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"

	// Seed the random number generator with the current time
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	random.Seed(time.Now().UnixNano())

	b := make([]byte, randStringLength)
	for i := range b {
		// randomly select 1 character from the given charset
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
