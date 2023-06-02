package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// IP Address = 169.6.7.20
// So, w = 169, x = 6, y = 7 and z = 20
// IP Number = 16777216*169 + 65536*6 + 256*7 + 20
//	= 2835349504 + 393216 + 1792 + 20
//	= 2835744532

func IPstringV4ToInt(str string) uint32 {
	octets := strings.Split(str, ".")
	if len(octets) != 4 {
		return 0
	}

	w, err := strconv.Atoi(octets[0])
	if err != nil {
		return 0
	}

	x, err := strconv.Atoi(octets[1])
	if err != nil {
		return 0
	}

	y, err := strconv.Atoi(octets[2])
	if err != nil {
		return 0
	}

	z, err := strconv.Atoi(octets[3])
	if err != nil {
		return 0
	}

	return uint32(256*256*256*w + 256*256*x + 256*y + z)
}

func IPIntV4ToString(n uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", n>>24&0xff, n>>16&0xff, n>>8&0xff, n&0xff)
}
