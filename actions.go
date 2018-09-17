package main

import (
	"fmt"

	// "github.com/spf13/viper"
	"github.com/WeMeetAgain/go-hdwallet"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"github.com/urfave/cli"
)

func seed2key(c *cli.Context) error {

	seedStr := c.String("seed")
	passPhrase := c.String("passPhrase")

	if seedStr == "" {
		return cli.NewExitError("No seed provided", 1)
	}

	seed := bip39.NewSeed(seedStr, passPhrase)

	// masterprv := hdwallet.MasterKey(seed)
	// masterpub := masterprv.Pub()

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	fmt.Println("Master private key: ", masterKey)
	fmt.Println("Master public key: ", publicKey)

	return nil
}

func derive(c *cli.Context) error {

	key := c.String("key")
	// path := c.String("path")

	if key == "" {
		return cli.NewExitError("No key provided", 1)
	}

	// Create a master private key
	masterprv, err := hdwallet.StringWallet(key)
	if err != nil {
		return err
	}

	// Convert a private key to public key
	masterpub := masterprv.Pub()

	// Generate new child key based on private or public key
	// childprv, err := masterprv.Child(0)
	childpub, err := masterpub.Child(0)

	// Create bitcoin address from public key
	address := childpub.Address()

	fmt.Println(address)

	return nil
}
