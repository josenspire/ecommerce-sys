package utils

import (
	"math/rand"
	"time"
)

const StrSrouce string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func GenerateRandString(length int) string {
	bytes := []byte(StrSrouce)
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func IsEmptyString(strArgs ...string) bool {
	var result = false
	for _, str := range strArgs {
		if str != "" {
			result = true
			break
		}
	}
	return result
}
