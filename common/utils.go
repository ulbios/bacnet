package common

import (
	"net"
	"strings"
)

// BoolToInt converts bool to int.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IntToBool converts int to bool.
func IntToBool(n int) bool {
	return n != 0
}

func IsLocalAddr(localAddrs []net.Addr, remoteAddr net.Addr) bool {
	for _, localAddr := range localAddrs {
		if strings.Split(localAddr.String(), "/")[0] == strings.Split(remoteAddr.String(), ":")[0] {
			return true
		}
	}
	return false
}
