/*
 * SPDX-License-Identifier: GPL-3.0-or-later
 *
 * Copyright (c) 2023. gsxab.
 *
 * This file is part of Software Update Watcher, a.k.a. Zhixin Robot.
 *
 * Software Update Watcher is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
 *
 * Software Update Watcher is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with Software Update Watcher. If not, see <https://www.gnu.org/licenses/>.
 */

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

func MatchExtractMultiple(matcher *regexp.Regexp, fullMatch bool, text string) (matched bool, results []string) {
	indices := matcher.FindStringSubmatchIndex(text)
	if len(indices) < 4 {
		return false, nil
	}
	if fullMatch && (indices[0] != 0 || indices[1] != len(text)) {
		return false, nil
	}
	results = make([]string, 0, len(indices)/2)
	for i := 2; i+1 < len(indices); i += 2 {
		results = append(results, text[indices[i]:indices[i+1]])
	}
	return true, results
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
