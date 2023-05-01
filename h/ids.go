package h

import "github.com/rs/xid"

func NewId(prefix string) string {
	guid := xid.New()
	return prefix + "_" + guid.String()
}

func NewIdPtr(prefix string) *string {
	value := NewId(prefix)
	return &value
}
