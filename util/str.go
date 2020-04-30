package util

import (
	"encoding/base64"
	"math/rand"
	"time"
)

var alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func MakeRandomAlphaString(n int) string {
	buf := make([]byte, n)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range buf {
		buf[i] = alphanum[int(rand.Int31n(int32(len(alphanum))))]
	}
	//return *(*string)(unsafe.Pointer(&b))
	return string(buf)
}

func MakeRandomBase64String(n int) string {
	b := make([]byte, n)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range b {
		b[i] = byte(rand.Int31n(255))
	}
	return base64.StdEncoding.EncodeToString(b)
}
