package util

import "regexp"

// MatchRegexGroups matches regexp and returns matching groups.
// returns false if not matched.
func MatchRegexGroups(r *regexp.Regexp, str string) ([]string, bool) {
	matches := r.FindStringSubmatch(str)
	if matches == nil {
		return []string{}, false
	}

	return matches[1:], true
}
