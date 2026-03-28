package settings

import "strings"

var userOverridablePrefixes = []string{
	"dashboard.",
	"navigation.",
}

func IsUserOverridable(key string) bool {
	for _, prefix := range userOverridablePrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}
