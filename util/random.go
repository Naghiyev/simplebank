package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomString(length int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwnerName() string {
	return RandomString(6)
}

func RandomEmailAddress() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomMoney() int {
	return (RandomInt(0, 1000))
}

func RandomCurrency() string {
	currencies := []string{"EUR", "GBP", "USD", "JPY"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
