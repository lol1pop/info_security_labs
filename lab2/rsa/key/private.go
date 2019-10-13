package key

import (
	"encoding/json"
	"github.com/DenisDyachkov/i_p_labs/basic"
	"io"
	"math/rand"
	"time"
)

type PrivateKey struct {
	N int64
	C int64
	f int64
}

func CreatePrivate() *PrivateKey {
	p, q := func() (_, _ int64) {
		generator := func() int64 {
			rand.Seed(time.Now().UnixNano())
			return rand.Int63n(20000) + 4 //todo: BigInt
		}
		q := basic.GeneratePrime(generator)
		p := basic.GeneratePrime(generator)
		for q == p {
			p = basic.GeneratePrime(generator)
		}
		return p, q
	}()
	n := p * q
	f := (p - 1) * (q - 1)
	c := calcPublicKey(f)
	return &PrivateKey{n, c, f}
}

func calcPublicKey(f int64) (c int64) {
	for basic.GCD(c, f) != 1 {
		c = rand.Int63n(f-4) + 3
	}
	return
}

func (k *PrivateKey) IsValid() bool {
	return k.N != 0 && k.C != 0
}

//todo: code gen (?)
func (k *PrivateKey) Save(wr io.Writer) error {
	bytes, err := json.Marshal(*k)
	if err == nil {
		_, err = wr.Write(bytes)
	}
	return err
}

func LoadPrivateKeyFromReader(rd io.Reader) (*PrivateKey, error) {
	var marshalled []byte
	if _, err := rd.Read(marshalled); err != nil {
		return nil, err
	}
	var key *PrivateKey
	err := json.Unmarshal(marshalled, key)
	return key, err
}
