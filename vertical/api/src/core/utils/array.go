package utils

import "slices"

func ContainsAtLeastOne(req []string, given []string) bool {
	for _, r := range req {
		if slices.Contains(given, r) {
			return true
		}
	}
	return false
}
