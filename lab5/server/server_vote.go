package server

import (
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa"
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa/key"
	"math/big"
)

func InitServer() *key.PublicKey {
	println("=== RSA ===")
	private, public := signature_rsa.CreateKeys()
	println("N ", private.N.String())
	println("C ", private.C.String())
	println("D ", public.D.String())
	return public
}

func ReadVote(vH *big.Int, key *key.PrivateKey) *big.Int {
	vS := new(big.Int).Exp(vH, key.C, key.N)
	return vS
}
