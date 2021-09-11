package mjwt

import (
	"time"
)

// Enum untuk tipe jwt
const (
	Access int = iota
	Refresh
)

type CustomClaim struct {
	Identity    int
	Name        string
	Exp         int64
	ExtraMinute time.Duration
	Type        int
	Fresh       bool
	Roles       []string
}
