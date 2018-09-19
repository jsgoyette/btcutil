package main

import (
	"github.com/urfave/cli"
)

var (
	hexFlag = cli.StringFlag{
		Name:  "hex",
		Value: "",
		Usage: "Generic hex string",
	}
	keyFlag = cli.StringFlag{
		Name:  "key, k",
		Value: "",
		Usage: "Private/Public key string (xprv/xpub)",
	}
	netFlag = cli.StringFlag{
		Name:  "net, n",
		Value: "",
		Usage: "Network (mainnet/testnet/regtest/simnet)",
	}
	pathFlag = cli.StringFlag{
		Name:  "path, p",
		Value: "0",
		Usage: "Derivation path string",
	}
	passPhraseFlag = cli.StringFlag{
		Name:  "passphrase, pass",
		Value: "",
		Usage: "Passphrase string",
	}
	seedFlag = cli.StringFlag{
		Name:  "seed, s",
		Value: "",
		Usage: "Seed words string",
	}

	commands = []cli.Command{
		{
			Name:        "keyFromSeed",
			Aliases:     []string{"key"},
			Usage:       "Get extended xprv key from seed",
			UsageText:   "btcutil seed2key [options]",
			Description: "Get the extended private key from a given seed.",
			Action:      keyFromSeed,
			Flags:       []cli.Flag{passPhraseFlag, seedFlag},
		},
		{
			Name:        "deriveFromKey",
			Aliases:     []string{"derive"},
			Usage:       "Derive child key",
			UsageText:   "btcutil derive [options]",
			Description: "Derive the extended private key from a given key.",
			Action:      derive,
			Flags:       []cli.Flag{keyFlag, pathFlag},
		},
		{
			Name: "hash160",
			// Aliases:     []string{"hash"},
			Usage:       "Generate hash160",
			UsageText:   "btcutil hash160 [options]",
			Description: "Generate the hash160 of a given hex string",
			Action:      hash160,
			Flags:       []cli.Flag{hexFlag},
		},
	}
)
