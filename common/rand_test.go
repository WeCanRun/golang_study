package common

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Perm(masters []string) {
	for i := len(masters); i > 0; i-- {
		lastIndex := i - 1
		randIndex := rand.Intn(i)
		masters[randIndex], masters[lastIndex] = masters[lastIndex], masters[randIndex]
	}
}

func TestPerm(t *testing.T) {
	var count = make(map[string]int)
	masters := []string{"10.124.142.220:5050", "10.124.142.221:5050", "10.124.142.222:5050"}
	for i := 0; i < 33333333; i++ {
		Perm(masters)
		count[masters[0]]++
	}
	t.Log("count: ", count)
}
