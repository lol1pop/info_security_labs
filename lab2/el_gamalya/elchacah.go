package el_gamalya

import (
	"encoding/binary"
	"encoding/json"
	"github.com/lol1pop/info_security_labs/basic"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
)

func choicePrimeNumbers() (_, _ uint64) {
	for {
		q := basic.RangeRandom(100000000, 10000000000)
		if basic.TestFerma(q, 15) {
			p := 2*q + 1
			if basic.TestFerma(p, 15) {
				return q, p
			}
		}
	}
}

func GeneratedPrimeNumbers() (_, _ uint64) {
	q, p := choicePrimeNumbers()
	for {
		max := int64(p) - 1
		g := uint64(rand.Int63n(max))
		if 1 < g && g < (p-1) {
			if basic.PowByModule(g, q, p) != 1 {
				return g, p
			}
		}
	}
}

type Key struct {
	P     uint64
	G     uint64
	OSkey uint64
}

func CreatedCoupleKeys(g, p uint64) (private, public Key) {
	secret := basic.RangeRandom(1, p-2)
	open := basic.PowByModule(g, secret, p)
	return Key{
			G:     g,
			P:     p,
			OSkey: secret,
		},
		Key{
			G:     g,
			P:     p,
			OSkey: open,
		}
}

func SecretSessionKey(p uint64) uint64 {
	return basic.RangeRandom(1, p)
}

func EncryptMessage(message []byte, publicKey Key) (rList []uint64, eList []uint64) {
	var encryptBuffer []uint64
	var rBuffer []uint64
	for _, m := range message {
		k := SecretSessionKey(publicKey.P)
		r := basic.PowByModule(publicKey.G, k, publicKey.P)
		c := encrypt(m, k, publicKey)
		encryptBuffer = append(encryptBuffer, c)
		rBuffer = append(rBuffer, r)
	}
	return rBuffer, encryptBuffer
}

func encrypt(m byte, k uint64, publicKey Key) uint64 {
	return (uint64(m) * basic.PowByModule(publicKey.OSkey, k, publicKey.P)) % publicKey.P
}

func decrypt(r, e uint64, privateKey Key) uint64 {
	power := privateKey.P - 1 - privateKey.OSkey
	return (e * basic.PowByModule(r, power, privateKey.P)) % privateKey.P
}

func DecryptMessage(r []uint64, e []uint64, privateKey Key) []byte {
	var decryptMessage []byte
	for i := 0; i < len(e); i++ {
		c := decrypt(r[i], e[i], privateKey)
		decryptMessage = append(decryptMessage, byte(c))
	}
	return decryptMessage
}

func EncryptMessageBinary(message []byte, publicKey Key) (r uint64, e []byte) {
	k := SecretSessionKey(publicKey.P)
	r = basic.PowByModule(publicKey.G, k, publicKey.P)
	var encryptBuffer []byte
	encrypted := make([]byte, 8)
	for _, m := range message {
		c := encrypt(m, k, publicKey)
		binary.LittleEndian.PutUint64(encrypted, c)
		encryptBuffer = append(encryptBuffer, encrypted...)
	}
	return r, encryptBuffer
}

func DecryptMessageBinary(r uint64, e []byte, key Key) []byte {
	var decryptMessage []byte
	for i := 0; i < len(e); i += 8 {
		g := binary.LittleEndian.Uint64(e[i : i+8])
		c := decrypt(r, g, key)
		decryptMessage = append(decryptMessage, byte(c))
	}
	return decryptMessage
}

func StartElGamalya() {
	g, p := GeneratedPrimeNumbers()
	private, public := CreatedCoupleKeys(g, p)
	msg := []byte{'e', 'l', ' ', 'c', 'h', 'a', 'c', 'h', 'a'}
	println(string(msg))
	r, e := EncryptMessage(msg, public)
	m := DecryptMessage(r, e, private)
	println(string(m))
}
