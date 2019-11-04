package signature_el_gam

import (
	"crypto/rand"
	"crypto/sha1"
	"math/big"
)

const NUM_BIT = 32

type PubKeys struct {
	Y *big.Int
	R *big.Int
	S *big.Int
}

type PubData struct {
	P   *big.Int
	G   *big.Int
	Key PubKeys

	Sm string
}

func initParams() (p, q, g *big.Int) {

	var err error
	for {
		q, err = rand.Prime(rand.Reader, NUM_BIT)
		if err != nil {
			panic(err)
		}
		p = new(big.Int).Mul(q, big.NewInt(2))
		p.Add(p, big.NewInt(1))
		if p.ProbablyPrime(50) {
			break
		}
	}

	max := new(big.Int).Sub(p, big.NewInt(1))

	for {
		g, err = rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		if g.Cmp(big.NewInt(1)) == 1 && g.Cmp(max) == -1 {
			if new(big.Int).Exp(g, q, p).Cmp(big.NewInt(1)) != 0 {
				break
			}
		}
	}

	return
}

func CreatedPrivateKey(p *big.Int) *big.Int {
	max := new(big.Int).Sub(p, big.NewInt(1))
	for {
		x, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		if x.Cmp(big.NewInt(1)) == 1 && x.Cmp(max) == -1 {
			return x
		}
	}
}

func CreatedPublicKey(p, g, x *big.Int) *big.Int {
	return new(big.Int).Exp(g, x, p)
}

func SignatureHash(m string) *big.Int {
	return new(big.Int).SetBytes(sha1.New().Sum([]byte(m)))
}

func SignatureElGamale(p, g, x *big.Int, m string) (r, s *big.Int) {
	h := SignatureHash(m)
	h = big.NewInt(4) //TODO()
	println("h ", h.Text(10))
	p_1 := new(big.Int).Sub(p, big.NewInt(1))
	var k *big.Int
	var err error

	for { //Проблема алгоритма в числе k  TODO()
		k, err = rand.Int(rand.Reader, p_1)
		if err != nil {
			panic(err)
		}
		if k.Cmp(big.NewInt(1)) == 1 && k.Cmp(p_1) == -1 {
			gcd := new(big.Int).GCD(nil, nil, k, p_1)
			if gcd.Cmp(big.NewInt(1)) == 0 {
				break
			}
		}
	}
	inverseK := new(big.Int).ModInverse(k, p_1)
	println("  -->", inverseK.Text(10))

	r = new(big.Int).Exp(g, k, p)
	xr := new(big.Int).Mul(x, r)
	u := h.Sub(h, xr).Mod(h, p_1)
	ku := new(big.Int).Mul(inverseK, u)
	s = new(big.Int).Mod(ku, p_1)

	return
}

func CheckSignatureElGamale(d PubData) bool {
	h := SignatureHash(d.Sm)
	h = big.NewInt(4) //TODO()
	yR := new(big.Int).Exp(d.Key.Y, d.Key.R, d.P)
	rS := new(big.Int).Exp(d.Key.R, d.Key.S, d.P)
	c1 := new(big.Int).Mul(yR, rS)
	c2 := new(big.Int).Exp(d.G, h, d.P)
	return c1.Cmp(c2) == 0
}

func StartElGamale() {
	m := "4"
	p, _, g := initParams()
	p = big.NewInt(23)
	g = big.NewInt(5)
	println("p ", p.Text(10))
	println("g ", g.Text(10))
	x := CreatedPrivateKey(p)
	y := CreatedPublicKey(p, g, x)
	x = big.NewInt(4)
	y = big.NewInt(4)
	println("x ", x.Text(10))
	println("y ", y.Text(10))
	r, s := SignatureElGamale(p, g, x, m)
	println("r ", r.Text(10))
	println("s ", s.Text(10))
	pubKeys := PubKeys{y, r, s}
	pubData := PubData{p, g, pubKeys, m}
	check := CheckSignatureElGamale(pubData)
	println(check)

}
