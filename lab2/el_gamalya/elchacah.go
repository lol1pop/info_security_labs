package el_gamalya

import (
	"github.com/lol1pop/info_security_labs/basic"
	"fmt"
	"math/rand"
)

func choicePrimeNumbers() (_,_ uint64) {
	for {
		q := basic.RangeRandom(50,1000000000)

		fmt.Print("|:")
		fmt.Println(q)
		if basic.TestFerma(q, 15) {
			fmt.Print("  -->: ")
			fmt.Println(q)
			p := 2 * q + 1
			if basic.TestFerma(p, 15){
				return q, p
			}
		}
	}
}

func GeneratedPrimeNumbers() (_,_ uint64) {
	q, p := choicePrimeNumbers()
	for {
		max := int64(p) - 1
		g := uint64(rand.Int63n(max))
		if 1 < g && g < (p - 1) {
			if basic.PowByModule(g, q, p) != 1 {
				return g, p
			}
		}
	}
}

type Key struct {
	p uint64
	g uint64
	key uint64
}

func CreatedCoupleKeys(g, p  uint64) (private, public Key) {
	secret := basic.RangeRandom(1,p)
	open := basic.PowByModule(g, secret, p)
	return Key{
		g: g,
		p: p,
		key: secret,
	},
	Key{
		g: g,
		p: p,
		key: open,
	}
}

func SecretSessionKey(p uint64) uint64 {
	return basic.RangeRandom(1, p)
}

func EncryptMessage(message []byte, publicKey Key) (r uint64, e []uint64) {
	k := SecretSessionKey(publicKey.p)
	r = basic.PowByModule(publicKey.g, k, publicKey.p)
	var encryptBuffer []uint64
	for _, m := range message {
		c := uint64(m) * basic.PowByModule(publicKey.key, k, publicKey.p)
		encryptBuffer = append(encryptBuffer, c)
	}
	return r, encryptBuffer
}

func decrypt(r, e uint64, privateKey Key) uint64 {
	power := privateKey.p - 1 - privateKey.key
	return e * basic.PowByModule(r, power, privateKey.p)
}

func DecryptMessage(r uint64, e []uint64, key Key) []byte {
	var decryptMessage []byte
	for i := 0; i < len(e); i++ {
		c := decrypt(r, e[i], key)
		decryptMessage = append(decryptMessage, byte(c))
	}
	return decryptMessage
}


func StartElGamalya(){
println("EL CHACHA")
	g,p := GeneratedPrimeNumbers()
	private, public := CreatedCoupleKeys(g,p)
	msg := []byte{'t'}
	r,e := EncryptMessage(msg, public)
	m := DecryptMessage(r,e, private)
	println(string(m))
}