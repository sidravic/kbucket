package util

import (
	"crypto/sha1"
	"time"
)

func GenerateRandomID() []byte {
	s := sha1.New()
	time := time.Now().String()
	s.Write([]byte(time))
	return s.Sum(nil)
}

func GenerateID(idString string) []byte {
	s := sha1.New()
	s.Write([]byte(idString))
	return s.Sum(nil)
}
