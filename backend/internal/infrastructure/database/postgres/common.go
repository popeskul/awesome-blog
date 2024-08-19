package postgres

import "strings"

func parseSortParam(sortParam string) (string, string) {
	sortParts := strings.Split(sortParam, "_")
	if len(sortParts) == 2 {
		field := sortParts[0]
		order := sortParts[1]
		if order == "asc" || order == "desc" {
			return field, order
		}
	}
	return "created_at", "desc"
}
