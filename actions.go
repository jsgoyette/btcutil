package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	// "github.com/spf13/viper"
	"github.com/tyler-smith/go-bip39"
	"github.com/urfave/cli"

	"github.com/btcsuite/btcd/txscript"
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
			publicKey.String(),
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
			derivedKey.String(),
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

func hash160(c *cli.Context) error {

	hexStr := c.String("hex")

	if hexStr == "" {
		return cli.NewExitError("No hex string provided", 1)
	}

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	hash := btcutil.Hash160(decoded)

	fmt.Printf("%s\n", hex.EncodeToString(hash))

	return nil
}
