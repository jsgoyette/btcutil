package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"hash"

	"fmt"

	// "github.com/spf13/viper"
	"github.com/tyler-smith/go-bip39"
	"github.com/urfave/cli"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"golang.org/x/crypto/ripemd160"
)

func keyFromSeed(c *cli.Context) error {

	seedStr := c.String("seed")
	passphrase := c.String("passphrase")

	if seedStr == "" {
		return cli.NewExitError("No seed provided", 1)
	}

	seed := bip39.NewSeed(seedStr, passphrase)

	masterKey, err := hdkeychain.NewMaster(seed, Network)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	publicKey, err := masterKey.Neuter()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	final := &Key{Network.Name, masterKey.String(), publicKey.String()}

	finalStr, err := json.Marshal(final)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("%s\n", finalStr)

	return nil
}

func derive(c *cli.Context) error {

	key := c.String("key")
	path := c.String("path")

	if key == "" {
		return cli.NewExitError("No key provided", 1)
	}

	masterKey, err := hdkeychain.NewKeyFromString(key)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	derivedKey := masterKey

	if path != "." && path != "m" {
		parsedPath, err := parsePath(path)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		for _, p := range parsedPath {
			derivedKey, err = derivedKey.Child(p)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}
	}

	publicKey, err := derivedKey.Neuter()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	x_prefix := PREFIX_xpub
	y_prefix := PREFIX_ypub
	z_prefix := PREFIX_zpub

	if Network.Name != "mainnet" {
		x_prefix = PREFIX_tpub
		y_prefix = PREFIX_upub
		z_prefix = PREFIX_vpub
	}

	x_version, err := hex.DecodeString(x_prefix)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	y_version, err := hex.DecodeString(y_prefix)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	z_version, err := hex.DecodeString(z_prefix)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	x_publicKey, err := NewKeyFromStringVersion(x_version, publicKey.String())
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	y_publicKey, err := NewKeyFromStringVersion(y_version, publicKey.String())
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	z_publicKey, err := NewKeyFromStringVersion(z_version, publicKey.String())
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	p2pkhAddress, err := derivedKey.Address(Network)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	ecPubKey, err := derivedKey.ECPubKey()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	keyHash := btcutil.Hash160(ecPubKey.SerializeCompressed())

	scriptSig, err := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(keyHash).Script()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	p2shAddress, err := btcutil.NewAddressScriptHash(scriptSig, Network)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	witnessProgram := btcutil.Hash160(ecPubKey.SerializeCompressed())

	bech32Address, err := btcutil.NewAddressWitnessPubKeyHash(witnessProgram, Network)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var finalStr []byte

	if masterKey.IsPrivate() {

		// add private key and wif
		ecPrivKey, err := derivedKey.ECPrivKey()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		wif, err := btcutil.NewWIF(ecPrivKey, Network, true)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		final := &DerivedPrivateKey{
			Network.Name,
			derivedKey.String(),
			x_publicKey.String(),
			y_publicKey.String(),
			z_publicKey.String(),
			wif.String(),
			hex.EncodeToString(ecPubKey.SerializeCompressed()),
			p2pkhAddress.String(),
			p2shAddress.String(),
			hex.EncodeToString(scriptSig),
			bech32Address.String(),
			hex.EncodeToString(witnessProgram),
		}

		finalStr, err = json.Marshal(final)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

	} else {

		final := &DerivedPublicKey{
			Network.Name,
			x_publicKey.String(),
			y_publicKey.String(),
			z_publicKey.String(),
			hex.EncodeToString(ecPubKey.SerializeCompressed()),
			p2pkhAddress.String(),
			p2shAddress.String(),
			hex.EncodeToString(scriptSig),
			bech32Address.String(),
			hex.EncodeToString(witnessProgram),
		}

		finalStr, err = json.Marshal(final)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	fmt.Printf("%s\n", finalStr)

	return nil
}

// btcsuite/btcutil/hash160.go
// Calculate the hash of hasher over buf.
func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

func calcRipemd160(c *cli.Context) error {

	hexStr := c.String("hex")

	if hexStr == "" {
		return cli.NewExitError("No hex string provided", 1)
	}

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	hashBytes := calcHash(decoded, ripemd160.New())

	fmt.Printf("%s\n", hex.EncodeToString(hashBytes))

	return nil
}

func calcHash160(c *cli.Context) error {

	hexStr := c.String("hex")

	if hexStr == "" {
		return cli.NewExitError("No hex string provided", 1)
	}

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	hashBytes := calcHash(calcHash(decoded, sha256.New()), ripemd160.New())

	fmt.Printf("%s\n", hex.EncodeToString(hashBytes))

	return nil
}

func calcHash256(c *cli.Context) error {

	hexStr := c.String("hex")

	if hexStr == "" {
		return cli.NewExitError("No hex string provided", 1)
	}

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	hashBytes := calcHash(calcHash(decoded, sha256.New()), sha256.New())

	fmt.Printf("%s\n", hex.EncodeToString(hashBytes))

	return nil
}

func calcSha256(c *cli.Context) error {

	hexStr := c.String("hex")

	if hexStr == "" {
		return cli.NewExitError("No hex string provided", 1)
	}

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	hashBytes := calcHash(decoded, sha256.New())

	fmt.Printf("%s\n", hex.EncodeToString(hashBytes))

	return nil
}
