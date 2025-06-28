package util

import "strings"

// FilterByPrefix returns a new slice with only the strings that start with the given prefix.
func FilterByPrefix(input []string, prefix string) []string {
	var filtered []string
	for _, s := range input {
		if strings.HasPrefix(s, prefix) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// FilterAndTrimPrefix returns a new slice containing only the strings
// that start with the given prefix, with the prefix removed.
func FilterAndTrimPrefix(input []string, prefix string) []string {
	var result []string
	for _, s := range input {
		if strings.HasPrefix(s, prefix) {
			result = append(result, strings.TrimPrefix(s, prefix))
		}
	}
	return result
}
