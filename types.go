package main

type Key struct {
	Network    string
	PrivateKey string
	PublicKey  string
}

type DerivedPrivateKey struct {
	Network        string
	PrivateKey     string
	PublicKey      string
	WIF            string
	ECPubKey       string
	P2PKHAddress   string
	P2SHAddress    string
	P2SHScript     string
	Bech32Address  string
	WitnessProgram string
}

type DerivedPublicKey struct {
	Network        string
	PublicKey      string
	ECPubKey       string
	P2PKHAddress   string
	P2SHAddress    string
	P2SHScript     string
	Bech32Address  string
	WitnessProgram string
}
