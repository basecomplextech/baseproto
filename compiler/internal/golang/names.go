// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import "strings"

func toUpperCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		part = strings.ToLower(part)
		part = strings.Title(part)
		parts[i] = part
	}

	s1 := strings.Join(parts, "")
	if strings.HasPrefix(s, "_") {
		s1 = "_" + s1
	}
	if strings.HasSuffix(s, "_") {
		s1 += "_"
	}
	return s1
}

func toLowerCameCase(s string) string {
	if len(s) == 0 {
		return ""
	}

	s = toUpperCamelCase(s)
	return strings.ToLower(s[:1]) + s[1:]
}
