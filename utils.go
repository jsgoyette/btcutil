package main

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
)

func parsePath(path string) ([]uint32, error) {
	parts := strings.Split(path, "/")
	nums := make([]uint32, len(parts), len(parts))

	for i, p := range parts {
		if len(p) == 0 {
			continue
		}

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

// NewKeyFromStringVersion modified from hdkeychain.NewKeyFromString
func NewKeyFromStringVersion(version []byte, key string) (*hdkeychain.ExtendedKey, error) {
	// The base58-decoded extended key must consist of a serialized payload
	// plus an additional 4 bytes for the checksum.
	decoded := base58.Decode(key)
	if len(decoded) != serializedKeyLen+4 {
		return nil, hdkeychain.ErrInvalidKeyLen
	}

	// The serialized format is:
	//   version (4) || depth (1) || parent fingerprint (4)) ||
	//   child num (4) || chain code (32) || key data (33) || checksum (4)

	// Split the payload and checksum up and ensure the checksum matches.
	payload := decoded[:len(decoded)-4]
	checkSum := decoded[len(decoded)-4:]
	expectedCheckSum := chainhash.DoubleHashB(payload)[:4]
	if !bytes.Equal(checkSum, expectedCheckSum) {
		return nil, hdkeychain.ErrBadChecksum
	}

	// Deserialize each of the payload fields.
	depth := payload[4:5][0]
	parentFP := payload[5:9]
	childNum := binary.BigEndian.Uint32(payload[9:13])
	chainCode := payload[13:45]
	keyData := payload[45:78]

	// The key data is a private key if it starts with 0x00.  Serialized
	// compressed pubkeys either start with 0x02 or 0x03.
	isPrivate := keyData[0] == 0x00
	if isPrivate {
		// Ensure the private key is valid.  It must be within the range
		// of the order of the secp256k1 curve and not be 0.
		keyData = keyData[1:]
		keyNum := new(big.Int).SetBytes(keyData)
		if keyNum.Cmp(btcec.S256().N) >= 0 || keyNum.Sign() == 0 {
			return nil, hdkeychain.ErrUnusableSeed
		}
	} else {
		// Ensure the public key parses correctly and is actually on the
		// secp256k1 curve.
		_, err := btcec.ParsePubKey(keyData, btcec.S256())
		if err != nil {
			return nil, err
		}
	}

	return hdkeychain.NewExtendedKey(version, keyData, chainCode, parentFP, depth,
		childNum, isPrivate), nil
}
