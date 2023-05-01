package h

import (
	"regexp"
	"strings"
)

var trailingSlashRegex = regexp.MustCompile("/+$")
var leadingSlashRegex = regexp.MustCompile("^/+")

func RemoveTrailingSlash(path string) string {
	return trailingSlashRegex.ReplaceAllString(strings.TrimSpace(path), "")
}

func RemoveLeadingSlash(path string) string {
	return leadingSlashRegex.ReplaceAllString(strings.TrimSpace(path), "")
}

func RemoveSlashes(path string) string {
	return RemoveLeadingSlash(RemoveTrailingSlash(path))
}

func NormalizeUri(path string) string {
	if path == "" || path == "/" {
		return ""
	}
	result := "/" + RemoveLeadingSlash(RemoveTrailingSlash(path))
	return result
}
