package main

const (
	PREFIX_xpub = "0488b21e"
	PREFIX_ypub = "049d7cb2"
	PREFIX_Ypub = "0295b43f"
	PREFIX_zpub = "04b24746"
	PREFIX_Zpub = "02aa7ed3"
	PREFIX_tpub = "043587cf"
	PREFIX_upub = "044a5262"
	PREFIX_Upub = "024289ef"
	PREFIX_vpub = "045f1cf6"
	PREFIX_Vpub = "02575483"

	serializedKeyLen = 4 + 1 + 4 + 4 + 32 + 33 // 78 bytes
)

type Key struct {
	Network    string
	PrivateKey string
	PublicKey  string
}

type DerivedPrivateKey struct {
	Network         string
	PrivateKey      string
	PublicKey       string
	PublicKeyP2SH   string
	PublicKeyP2WPKH string
	WIF             string
	ECPubKey        string
	P2PKHAddress    string
	P2SHAddress     string
	P2SHScript      string
	Bech32Address   string
	WitnessProgram  string
}

type DerivedPublicKey struct {
	Network         string
	PublicKey       string
	PublicKeyP2SH   string
	PublicKeyP2WPKH string
	ECPubKey        string
	P2PKHAddress    string
	P2SHAddress     string
	P2SHScript      string
	Bech32Address   string
	WitnessProgram  string
}
