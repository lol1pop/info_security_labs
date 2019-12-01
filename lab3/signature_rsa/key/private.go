package key

import (
	"crypto/rand"
	basic "github.com/DenisDyachkov/i_p_labs/basic/big"
	"math/big"
)

const RSAbits int = 32

type PrivateKey struct {
	N *big.Int `json:"n"`
	C *big.Int `json:"c"`
	f *big.Int
}

func CreatePrivate() *PrivateKey {
	p, q, err := func() (*big.Int, *big.Int, error) {
		q, err := rand.Prime(rand.Reader, RSAbits)
		if err != nil {
			return nil, nil, err
		}
		p, err := rand.Prime(rand.Reader, RSAbits)
		if err != nil {
			return nil, nil, err
		}
		for q.Cmp(p) == 0 {
			p, _ = rand.Prime(rand.Reader, RSAbits)
		}
		return p, q, nil
	}()
	if err != nil {
		return nil
	}
	n := new(big.Int).Mul(p, q)
	f := new(big.Int).Mul(new(big.Int).Sub(p, basic.One), new(big.Int).Sub(q, basic.One))
	c := calcPublicKey(f)
	return &PrivateKey{n, c, f}
}

func calcPublicKey(f *big.Int) *big.Int {
	max := new(big.Int).Sub(f, basic.One)
	c, _ := rand.Int(rand.Reader, max)
	for basic.GCD(c, f).Cmp(basic.One) != 0 {
		c, _ = rand.Int(rand.Reader, max)
	}
	return c
}
