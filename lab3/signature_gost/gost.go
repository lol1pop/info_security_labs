package signature_gost

import (
	"crypto/rand"
	"crypto/sha1"
	"math/big"
)

const Q_BIT = 15
const B_BIT = 16
const P_BIT = 30

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

func find_p(q *big.Int) (_, _ *big.Int) {
	for {
		p, err := rand.Prime(rand.Reader, P_BIT)
		if err != nil {
			panic(err)
		}
		b := new(big.Int)
		b.Sub(p, big.NewInt(1))
		if new(big.Int).Mod(b, q).Cmp(big.NewInt(0)) == 0 {
			b.Div(b, q)
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
	p, b := find_p(q)
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

func CheckSignatureGost(r, s *big.Int, m string, d ExchangeData) bool {
	h := new(big.Int).SetBytes(sha1.New().Sum([]byte(m)))
	iH := new(big.Int).ModInverse(h, d.Data.Q) //??????????
	println(iH.Text(10))

	u1 := new(big.Int).Mod(new(big.Int).Mul(s, iH), d.Data.Q)
	u2 := new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Neg(r), iH), d.Data.Q)
	au1 := new(big.Int).Exp(d.Data.A, u1, d.Data.P)
	yu2 := new(big.Int).Exp(d.PublicKey, u2, d.Data.P)
	v := new(big.Int).Mod(new(big.Int).Mod(new(big.Int).Mul(au1, yu2), d.Data.P), d.Data.Q)
	return v.Cmp(r) == 0
}

func StartGost() {

	pubData := InitParams()
	println(" ======= \nQ ", pubData.Q.Text(10))
	println("P ", pubData.P.Text(10))
	println("A ", pubData.A.Text(10), " \n ======= ")
	m := "hello"
	pubA, pivA := CreatePublicPrivateKey(pubData)
	println("k1 ", pubA.Text(10))
	println("K2 ", pivA.Text(10))
	r, s := SignatureGost(pubData, pivA, m)
	println("r ", r.Text(10))
	println("s ", s.Text(10))
	data := ExchangeData{pubData, pubA}
	check := CheckSignatureGost(r, s, m, data)
	println(check)

}
