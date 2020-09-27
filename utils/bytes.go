package utils

import (
	"fmt"
)

// DebugBytes prints out the hex equivalent
func DebugBytes(buf []byte) string {
	debug := ""
	for i := 0; i < len(buf); i++ {
		debug += fmt.Sprintf("0x%02x ", buf[i])
	}
	return debug
}
