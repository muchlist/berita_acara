package dao

import "strings"

const (
	ErrSqlBuilder = "kesalahan pada sql builder"
	ErrRollback   = "gagal ketika roleback"
	ErrCommit     = "gagal ketika commit"
)

func Returning(columns ...string) string {
	sb := strings.Builder{}
	sb.WriteString("RETURNING ")
	for _, key := range columns {
		sb.WriteString(key + ", ")
	}
	return strings.TrimSuffix(sb.String(), ", ")
}
