package provider

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

// stringChecksum takes a string and returns the checksum of the string.
func stringChecksum(s string) string {
	h := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", h)
}

// stringListChecksum takes a slice of strings and returns the checksum of the strings.
func stringListChecksum(s []string) string {
	sort.Strings(s)
	return stringChecksum(strings.Join(s, ""))
}
