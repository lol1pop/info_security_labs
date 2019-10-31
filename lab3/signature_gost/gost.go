package signature_gost

import (
	"crypto/rand"
	"crypto/sha1"
	"github.com/lol1pop/info_security_labs/basic"
	"math/big"
)

type PublicData struct {
	P *big.Int
	Q *big.Int
	A *big.Int
}

type ExchangeData struct {
	Data      PublicData
	PublicKey *big.Int
}

func create_param_a(b, q, p *big.Int) (a *big.Int) {
	for {
		var g *big.Int
		for {
			g, _ = rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(1)))
			if g.ProbablyPrime(50) {
				break
			}
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
		if p.BitLen() == 1024 && p.ProbablyPrime(50) {
			return p, b
		}
	}
}

func CreatePublicPrivateKey(d PublicData) (_public, _private *big.Int) {
	x, err := rand.Int(rand.Reader, d.Q)
	if err != nil {
		panic(err)
	}
	y := basic.BigFastPowByModule(d.A, x, d.P)
	return y, x
}

func InitParams() PublicData {
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
	a := create_param_a(b, q, p)
	return PublicData{p, q, a}
}

func SignatureGost(d PublicData, x *big.Int, m string) (_, _ *big.Int) {
	//h =H(m) + [0 < h < q]
	var k, r, s, h *big.Int
	h = new(big.Int).SetBytes(sha1.New().Sum([]byte(m)))
	for {
		k, _ = rand.Int(rand.Reader, d.Q)
		r = basic.BigFastPowByModule(d.A, k, d.P)
		r.Mod(r, d.Q)
		if r.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		kh := new(big.Int).Mul(k, h)
		xr := new(big.Int).Mul(x, r)
		s = new(big.Int).Add(kh, xr)
		s.Mod(s, d.Q)
		if s.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		if d.Q.Cmp(r) != 1 && d.Q.Cmp(s) != 1 {
			continue
		}
		return r, s
	}
}

func ChecSignatureGost(r, s *big.Int, m string, d ExchangeData) bool {
	h := new(big.Int).SetBytes(sha1.New().Sum([]byte(m)))
	invers_h := new(big.Int).ModInverse(h, d.Data.P) //??????????
	sh := new(big.Int).Mul(s, invers_h)
	rh := new(big.Int).Mul(r, invers_h)
	u1 := new(big.Int).Mod(sh, d.Data.Q)
	u2 := new(big.Int).Mod(rh, d.Data.Q)
	u2.Neg(u2)
	au1 := basic.BigFastPowByModule(d.Data.A, u1, d.Data.P)
	yu2 := basic.BigFastPowByModule(d.PublicKey, u1, d.Data.P)
	v := new(big.Int).Mul(au1, yu2)
	v.Mod(v, d.Data.Q)
	return v.Cmp(r) == 0
}

func StartGost() {

	println(basic.BigFastPowByModule(big.NewInt(5), big.NewInt(7), big.NewInt(47)).Text(10))

	pubData := InitParams()
	println(pubData.P.Text(10))
	println(pubData.Q.Text(10))
	println(pubData.A.Text(10))
	/* = Generated = */
	/* == Alisa == */
	pubA, pivA := CreatePublicPrivateKey(pubData)
	println(pubA.Text(10), pivA.Text(10))
	/* ==  Bob  == */

}
