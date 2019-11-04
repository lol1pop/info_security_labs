package signature_gost

import (
	"crypto/rand"
	"crypto/sha1"
	"math/big"
)

const Q_BIT = 15
const B_BIT = 16
const P_BIT = 1024

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
		g, _ := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(1)))
		a = new(big.Int).Exp(g, b, p)
		if a.Cmp(big.NewInt(1)) == 1 {
			return a
		}
	}
}

func find_p_via_rand(q *big.Int) (_, _ *big.Int) {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(B_BIT), nil)
	for {
		b, _ := rand.Int(rand.Reader, max)
		p := b.Mul(b, q).Add(b, big.NewInt(1))
		if p.ProbablyPrime(50) {
			return p, b
		}
	}
}

func CreatePublicPrivateKey(d PublicData) (_public, _private *big.Int) {
	x, err := rand.Int(rand.Reader, d.Q)
	if err != nil {
		panic(err)
	}
	y := new(big.Int).Exp(d.A, x, d.P)
	return y, x
}

func InitParams() PublicData {
	q, err := rand.Prime(rand.Reader, Q_BIT)
	if err != nil {
		panic(err)
	}
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
		r = new(big.Int).Exp(d.A, k, d.P)
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
	invers_h := new(big.Int).ModInverse(h, d.Data.Q) //??????????
	println(invers_h.Text(10))

	sh := new(big.Int).Mul(s, invers_h)
	rh := new(big.Int).Mul(r, invers_h)
	rh.Neg(rh)
	u1 := new(big.Int).Mod(sh, d.Data.Q)
	u2 := new(big.Int).Mod(rh, d.Data.Q)
	au1 := new(big.Int).Exp(d.Data.A, u1, d.Data.P)
	yu2 := new(big.Int).Exp(d.PublicKey, u2, d.Data.P)
	v := new(big.Int).Mul(au1, yu2)
	v.Mod(v, d.Data.P)
	v.Mod(v, d.Data.Q)
	return v.Cmp(r) == 0
}

func SignatureGostTEST(d PublicData, x *big.Int) (r, s *big.Int) {
	//h =H(m) + [0 < h < q]
	var k, h *big.Int
	h = big.NewInt(4)
	for {
		k, _ = rand.Int(rand.Reader, d.Q)
		r = new(big.Int).Exp(d.A, k, d.P)
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

func ChecSignatureGostTEST(r, s *big.Int, d ExchangeData) bool {
	h := big.NewInt(4)
	invers_h := new(big.Int).ModInverse(h, d.Data.Q) //??????????
	println(invers_h.Text(10))

	sh := new(big.Int).Mul(s, invers_h)
	rh := new(big.Int).Mul(r, invers_h)
	rh.Neg(rh)
	u1 := new(big.Int).Mod(sh, d.Data.Q)
	u2 := new(big.Int).Mod(rh, d.Data.Q)
	au1 := new(big.Int).Exp(d.Data.A, u1, d.Data.P)
	yu2 := new(big.Int).Exp(d.PublicKey, u2, d.Data.P)
	v := new(big.Int).Mul(au1, yu2)
	v.Mod(v, d.Data.P)
	v.Mod(v, d.Data.Q)
	return v.Cmp(r) == 0
}

func StartGost() {

	for {
		pubData := InitParams()

		//p := big.NewInt(2111)
		//q := big.NewInt(211)
		//a := big.NewInt(196)
		////a = create_param_a(big.NewInt(22), q, p)
		//pubData = PublicData{p, q, a}

		println(" ======= \n", pubData.Q.Text(10), "   ", pubData.Q.BitLen())
		println(pubData.P.Text(10), "   ", pubData.P.BitLen())
		println(pubData.A.Text(10), "   ", pubData.P.BitLen(), " \n ======= ")
		//2994
		//42
		pubA, pivA := CreatePublicPrivateKey(pubData)
		println(pubA.Text(10))
		println(pivA.Text(10))
		r, s := SignatureGostTEST(pubData, pivA)
		println(r.Text(10))
		println(s.Text(10))
		lol := ExchangeData{pubData, pubA}
		check := ChecSignatureGostTEST(r, s, lol)
		println(check)
		println(check)
	}

	///* = Generated = */
	//pubData := InitParams()
	//println(pubData.P.Text(10))
	//println(pubData.Q.Text(10))
	//println(pubData.A.Text(10))
	///* == Alisa == */
	//pubA, pivA := CreatePublicPrivateKey(pubData)
	//println(pubA.Text(10),"\n", pivA.Text(10))
	//
	//data := ExchangeData{ pubData, pubA}
	//m := "hello"
	//r,s := SignatureGost(pubData, pivA, m)
	/* ==  Bob  == */
	//check := ChecSignatureGost(r,s,m, data)
	//println(check)

}
