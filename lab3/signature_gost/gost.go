package signature_gost

import (
	"crypto/rand"
	"github.com/lol1pop/info_security_labs/basic"
	"math/big"
)

func create_param_a(b, q, p *big.Int) (a *big.Int) {
	for {
		_, err := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(1)))
		g, _ := rand.Prime(rand.Reader, 1023)
		if err != nil {
			panic(err)
		}
		a = basic.BigFastPowByModule(g, b, p)
		if a.Cmp(big.NewInt(1)) == 1 {
			check_params_a := basic.BigFastPowByModule(a, q, p)
			println(check_params_a.Text(10))
			if check_params_a.Cmp(big.NewInt(1)) == 0 {
				return a
			}
		}
	}
}

func find_p_via_rand(q *big.Int) (_, _ *big.Int) {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(769), nil)
	for {
		b, _ := rand.Int(rand.Reader, max)
		p := b.Mul(b, q).Add(b, big.NewInt(1))
		if p.BitLen() == 1024 {
			if p.ProbablyPrime(50) {
				return p, b
			}
		}
	}
}

func CreatePublicPrivateKey(a, q, p *big.Int) (_public, _private *big.Int) {
	x, err := rand.Int(rand.Reader, q)
	if err != nil {
		panic(err)
	}
	y := basic.BigFastPowByModule(a, x, p)
	return y, x
}

func InitParams() (p, q, a *big.Int) {
	q, err := rand.Prime(rand.Reader, 256)
	if err != nil {
		panic(err)
	}
	//p, err = rand.Prime(rand.Reader, 1024)
	//if err != nil {
	//	panic(err)
	//}
	//b := new(big.Int)
	//b.Sub(p, big.NewInt(1))
	//b.Div(b, q)
	p, b := find_p_via_rand(q)
	a = create_param_a(b, q, p)
	return p, q, a
}

func StartGost() {

	println(basic.BigFastPowByModule(big.NewInt(5), big.NewInt(7), big.NewInt(47)).Text(10))

	p, q, a := InitParams()
	println(p.Text(10))
	println(q.Text(10))
	println(a.Text(10))
	/* = Generated = */
	/* == Alisa == */
	pubA, pivA := CreatePublicPrivateKey(a, q, p)
	println(pubA.Text(10), pivA.Text(10))
	/* ==  Bob  == */

}
