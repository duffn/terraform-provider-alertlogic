package provider

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

// stringChecksum takes a string and returns the checksum of the string.
func stringChecksum(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

// stringListChecksum takes a slice of strings and returns the checksum of the strings.
func stringListChecksum(s []string) string {
	sort.Strings(s)
	return stringChecksum(strings.Join(s, ""))
}
