package key

import (
	"encoding/json"
	"github.com/DenisDyachkov/i_p_labs/basic"
	"io"
)

type PublicKey struct {
	N int64
	D int64
}

func CreatePublic(privateKey *PrivateKey) *PublicKey {
	d, _ := basic.InversionByModule(privateKey.C, privateKey.f)
	return &PublicKey{privateKey.N, d}
}

func (k *PublicKey) IsValid() bool {
	return k.N != 0 && k.D != 0
}

func (k *PublicKey) Save(wr io.Writer) error {
	bytes, err := json.Marshal(*k)
	if err == nil {
		_, err = wr.Write(bytes)
	}
	return err
}

func LoadPublicKeyFromReader(rd io.Reader) (*PublicKey, error) {
	var marshalled []byte
	if _, err := rd.Read(marshalled); err != nil {
		return nil, err
	}
	var key *PublicKey
	err := json.Unmarshal(marshalled, key)
	return key, err
}
