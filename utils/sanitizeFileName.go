package utils

import "regexp"

func SanitizeFileName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
	name = re.ReplaceAllString(name, "_")

	if len(name) > 50 {
		name = name[:50]
	}
	return name
}
