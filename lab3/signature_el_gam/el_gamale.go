package signature_el_gam

import (
	"crypto/rand"
	"encoding/json"
	"github.com/lol1pop/info_security_labs/basic"
	"io/ioutil"
	"math/big"
	"os"
)

const NUM_BIT = 256

type PubKeys struct {
	Y *big.Int `json:"y"`
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

type PubData struct {
	P   *big.Int `json:"p"`
	G   *big.Int `json:"g"`
	Key PubKeys  `json:"key"`
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

func SignatureElGamale(p, g, x *big.Int, m []byte) (r, s *big.Int) {
	h := basic.GetMessageHash(m)
	println("h ", h.Text(10))
	p_1 := new(big.Int).Sub(p, big.NewInt(1))
	var k *big.Int
	var err error

	for {
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

	r = new(big.Int).Exp(g, k, p)
	xr := new(big.Int).Mul(x, r)
	u := h.Sub(h, xr).Mod(h, p_1)
	ku := new(big.Int).Mul(inverseK, u)
	s = new(big.Int).Mod(ku, p_1)

	return
}

func CheckSignatureElGamale(d PubData, m []byte) bool {
	h := basic.GetMessageHash(m)
	yR := new(big.Int).Exp(d.Key.Y, d.Key.R, d.P)
	rS := new(big.Int).Exp(d.Key.R, d.Key.S, d.P)
	c1 := new(big.Int).Mod(new(big.Int).Mul(yR, rS), d.P)
	c2 := new(big.Int).Exp(d.G, h, d.P)
	return c1.Cmp(c2) == 0
}

func StartElGamale() {
	m := "test message el Cha chA"
	_ = ioutil.WriteFile("ElGamaleSign.txt", []byte(m), os.ModePerm)
	fsrc, _ := ioutil.ReadFile("ElGamaleSign.txt")
	println(string(fsrc))
	p, _, g := initParams()
	println("p ", p.Text(10))
	println("g ", g.Text(10))
	x := CreatedPrivateKey(p)
	y := CreatedPublicKey(p, g, x)
	println("x ", x.Text(10))
	println("y ", y.Text(10))
	r, s := SignatureElGamale(p, g, x, fsrc)
	println("r ", r.Text(10))
	println("s ", s.Text(10))
	pubKeys := PubKeys{y, r, s}
	pubData := PubData{p, g, pubKeys}
	toJson := func(any interface{}) []byte {
		bytes, _ := json.Marshal(any)
		return bytes
	}
	_ = ioutil.WriteFile("ElGamaleSign-sing.txt", toJson(pubData), os.ModePerm)
	println("Check result: ", CheckSignatureElGamale(pubData, fsrc))
	CheckSignFile()
}

func CheckSignFile() {
	fsrc, _ := ioutil.ReadFile("ElGamaleSign.txt")
	src, _ := ioutil.ReadFile("ElGamaleSign-sing.txt")
	var data PubData
	_ = json.Unmarshal(src, &data)
	println("Check result: ", CheckSignatureElGamale(data, fsrc))
}
