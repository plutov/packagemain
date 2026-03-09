package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

func getCounter(t time.Time) uint64 {
	return uint64(t.Unix()) / 30
}

func generateHMAC(secret string, counter uint64) ([]byte, error) {
	secret = strings.TrimRight(strings.ToUpper(secret), "=")
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return nil, err
	}

	// The uint64 counter must be encoded as an 8-byte big-endian byte array
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	h := hmac.New(sha1.New, key)
	h.Write(buf)
	return h.Sum(nil), nil
}

func truncate(hash []byte) int {
	offset := hash[len(hash)-1] & 0x0f

	code := int(hash[offset]&0x7f)<<24 |
		int(hash[offset+1])<<16 |
		int(hash[offset+2])<<8 |
		int(hash[offset+3])

	return code % 1000000
}

func GenerateTOTP(secret string) (string, int, error) {
	counter := getCounter(time.Now())
	hash, err := generateHMAC(secret, counter)
	if err != nil {
		return "", 0, err
	}

	code := truncate(hash)

	timeRemaining := 30 - (int(time.Now().Unix()) % 30)

	return fmt.Sprintf("%d", code), timeRemaining, nil
}
