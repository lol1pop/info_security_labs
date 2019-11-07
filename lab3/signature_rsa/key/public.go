package key

import (
	"math/big"
)

type PublicKey struct {
	N *big.Int `json:"n"`
	D *big.Int `json:"d"`
}

func CreatePublic(privateKey *PrivateKey) *PublicKey {
	d := new(big.Int).ModInverse(privateKey.C, privateKey.f)
	return &PublicKey{privateKey.N, d}
}
