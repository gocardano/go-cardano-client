package utils

import (
	"time"
)

// TimeNowLower32 returns the lower 32 bits of the sender's monotonic clock
func TimeNowLower32() uint32 {
	return uint32(time.Now().UnixNano() & 0xffffffff)
}
