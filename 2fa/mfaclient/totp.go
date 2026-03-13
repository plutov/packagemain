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

func calculateHash(secret string, counter uint64) ([]byte, error) {
	secret = strings.ReplaceAll(secret, " ", "")
	secret = strings.TrimRight(secret, "=")
	secret = strings.ToUpper(secret)

	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return nil, err
	}

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

func GenerateTotp(secret string) (string, int, error) {
	now := time.Now()

	counter := getCounter(now)

	hash, err := calculateHash(secret, counter)
	if err != nil {
		return "", 0, err
	}

	code := truncate(hash)

	remaining := 30 - (int(now.Unix()) % 30)

	return fmt.Sprintf("%d", code), remaining, nil
}
