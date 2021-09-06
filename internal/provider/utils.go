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

// contains checks if a string is present in a slice.
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// expandInterfaceToStringList turns an interface into a string slice.
func expandInterfaceToStringList(list interface{}) []string {
	ifaceList := list.([]interface{})
	vs := make([]string, 0, len(ifaceList))
	for _, v := range ifaceList {
		vs = append(vs, v.(string))
	}
	return vs
}
