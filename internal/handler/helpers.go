package handler

import (
	"strconv"
	"strings"
)

// parseIDFromPath extracts numeric ID from URL path
// Example: "/api/categories/5" with prefix "/api/categories/" returns 5
func parseIDFromPath(path, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)
	return strconv.Atoi(idStr)
}
