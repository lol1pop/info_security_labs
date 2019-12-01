package course

import (
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa/key"
	"math/big"
)

type Node struct {
	Edge       int
	Color      int
	R          *big.Int
	PrivateKey *key.PrivateKey
	PublicKey  *key.PublicKey
}

type Edge struct {
	Weight int
	Start  Node
	End    Node
}

type EncryptNode struct {
	N    *big.Int
	D    *big.Int
	Z    *big.Int
	Edge int
}
