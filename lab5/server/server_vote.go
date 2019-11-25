package server

import (
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa"
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa/key"
	"github.com/lol1pop/info_security_labs/lab5/client"
	"math/big"
)

type Server struct {
	privateKey *key.PrivateKey
	PublicKey  *key.PublicKey
	voterSlice []*big.Int
}

func (s Server) toPrintKeys() {
	println("N ", s.privateKey.N.String())
	println("C ", s.privateKey.C.String())
	println("D ", s.PublicKey.D.String())
}

func (s *Server) InitServer() *Server {
	println("=== RSA ===")
	private, public := signature_rsa.CreateKeys()
	s.privateKey, s.PublicKey = private, public
	s.voterSlice = []*big.Int{}
	s.toPrintKeys()
	return s
}

func (s Server) ReadVote(vH *big.Int) *big.Int {
	vS := new(big.Int).Exp(vH, s.privateKey.C, s.privateKey.N)
	return vS
}

func (s Server) NewClient() client.Client {
	return client.Client{PublicKey: s.PublicKey}
}
