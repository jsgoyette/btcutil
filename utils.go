package main

import (
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/hdkeychain"
)

func parsePath(path string) ([]uint32, error) {
	parts := strings.Split(path, "/")
	nums := make([]uint32, len(parts), len(parts))

	for i, p := range parts {
		harden := uint32(0)
		lastChar := p[len(p)-1:]

		if lastChar == "'" || lastChar == "h" {
			p = p[:len(p)-1]
			harden = hdkeychain.HardenedKeyStart
		}

		num, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			return nil, err
		}

		nums[i] = uint32(num) + harden
	}

	return nums, nil
}
