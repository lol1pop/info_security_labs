package course

import (
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa/key"
	"math/big"
)

type Node struct {
	Edge       int             `json:"edge"`
	Color      int             `json:"color"`
	R          *big.Int        `json:"r,omitempty"`
	PrivateKey *key.PrivateKey `json:"private_key,omitempty"`
	PublicKey  *key.PublicKey  `json:"public_key,omitempty"`
}

type Edge struct {
	Weight int
	Start  Node
	End    Node
}

type EncryptNode struct {
	N    *big.Int `json:"n"`
	D    *big.Int `json:"d"`
	Z    *big.Int `json:"z"`
	Edge int      `json:"edge"`
}
