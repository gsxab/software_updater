package util

import "regexp"

func MatchExtract(matcher *regexp.Regexp, fullMatch bool, text string) (matched bool, result string) {
	indices := matcher.FindStringSubmatchIndex(text)
	if len(indices) < 4 {
		return false, ""
	}
	if fullMatch && (indices[0] != 0 || indices[1] != len(text)) {
		return false, ""
	}
	return true, text[indices[2]:indices[3]]
}

func Match(matcher *regexp.Regexp, fullMatch bool, text string) (matched bool) {
	indices := matcher.FindStringIndex(text)
	if len(indices) < 2 {
		return false
	}
	if fullMatch && (indices[0] != 0 || indices[1] != len(text)) {
		return false
	}
	return true
}
