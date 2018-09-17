package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	// "github.com/spf13/viper"
	"github.com/tyler-smith/go-bip39"
	"github.com/urfave/cli"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
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

	// masterKey, _ := bip32.NewMasterKey(seed)
	// publicKey := masterKey.PublicKey()

	finalKey := &Key{masterKey.String(), publicKey.String(), "", "", ""}

	final, err := json.Marshal(finalKey)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("%s\n", final)

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

	parsedPath, err := parsePath(path)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	derivedKey := masterKey

	for _, p := range parsedPath {
		derivedKey, err = derivedKey.Child(p)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	publicKey, err := derivedKey.Neuter()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	address, err := derivedKey.Address(Network)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	ecPrivKey, err := derivedKey.ECPrivKey()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	wif, err := btcutil.NewWIF(ecPrivKey, Network, true)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	ecPubKey, err := derivedKey.ECPubKey()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	finalKey := &Key{
		derivedKey.String(),
		publicKey.String(),
		wif.String(),
		hex.EncodeToString(ecPubKey.SerializeCompressed()),
		address.String(),
	}

	final, err := json.Marshal(finalKey)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("%s\n", final)

	return nil
}
