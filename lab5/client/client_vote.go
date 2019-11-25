package client

import (
	"crypto/rand"
	"github.com/lol1pop/info_security_labs/basic"
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa/key"
	"math/big"
)

type Client struct {
	PublicKey *key.PublicKey
}

type Voter struct {
	Name   string
	Vote   int64
	VH     *big.Int
	h      *big.Int
	rnd    *big.Int
	n      *big.Int
	r      *big.Int
	client *Client
}

func (c *Client) InitVote(name string, vote int64) Voter {
	max := new(big.Int).Exp(big.NewInt(2), big.NewInt(512), nil)
	rnd, _ := rand.Int(rand.Reader, max)
	v := big.NewInt(vote)
	n, _ := new(big.Int).SetString(rnd.String()+v.String(), 10)
	r := func(n *big.Int) *big.Int {
		for {
			max := new(big.Int).Exp(big.NewInt(2), big.NewInt(512), nil)
			r, _ := rand.Int(rand.Reader, max)
			gcd := new(big.Int).GCD(nil, nil, r, n)
			if gcd.Cmp(big.NewInt(1)) == 0 {
				return r
			}
		}
	}(c.PublicKey.N)

	h := basic.GetMessageHash([]byte(n.String()))
	println("n  " + n.String())
	println("h  " + h.String())

	vH := func(h, r *big.Int) *big.Int {
		return new(big.Int).Mod(new(big.Int).Mul(h, new(big.Int).Exp(r, c.PublicKey.D, c.PublicKey.N)), c.PublicKey.N)
	}(h, r)

	return Voter{name, vote, vH, h, rnd, n, r, c}
}

func (v Voter) Signature(vS *big.Int) (_, _ *big.Int) {
	iR := new(big.Int).ModInverse(v.r, v.client.PublicKey.N)
	s := new(big.Int).Mod(new(big.Int).Mul(vS, iR), v.client.PublicKey.N)
	return v.n, s
}
