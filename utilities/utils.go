package utilities

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

func Normalize(input string) string {
	// Step 1: Replace non-ascii or special characters with "-"
	re := regexp.MustCompile(`[^a-zA-Z0-9-._]+`)
	input = re.ReplaceAllString(input, "-")

	// Step 2: Replace consecutive occurrences of `.`, `-`, `_` with a single `-`
	re = regexp.MustCompile(`[.]{2,}`)
	input = re.ReplaceAllString(input, "-")

	re = regexp.MustCompile(`[-]{2,}`)
	input = re.ReplaceAllString(input, "-")

	re = regexp.MustCompile(`[_]{2,}`)
	input = re.ReplaceAllString(input, "-")

	return strings.ToLower(input)
}

// Hash returns the hex encoded sha256 hash of the given data
func Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
