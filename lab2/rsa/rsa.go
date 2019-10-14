package RSA

import (
	"encoding/binary"
	"github.com/DenisDyachkov/i_p_labs/basic"
	"github.com/DenisDyachkov/i_p_labs/crypt/RSA/key"
	"io/ioutil"
	"os"
)

func CreateKeys() (*key.PrivateKey, *key.PublicKey, error) {
	save := func(key key.Key, filename string) error {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		return key.Save(file)
	}
	privateKey := key.CreatePrivate()
	publicKey := key.CreatePublic(privateKey)
	_ = save(privateKey, "private.rsa")
	_ = save(publicKey, "public.rsa")
	return privateKey, publicKey, nil
}

func Encrypt(m int64, pub *key.PublicKey) int64 {
	if m >= pub.N {
		panic("m >= N")
	}
	return basic.FastPowByModule(m, pub.D, pub.N)
}

func Decrypt(e int64, priv *key.PrivateKey) int64 {
	return basic.FastPowByModule(e, priv.C, priv.N)
}

func EncryptBuffer(msg []byte, pub *key.PublicKey) []byte {
	var result []byte
	encrypted := make([]byte, 8)
	for _, m := range msg {
		c := Encrypt(int64(m), pub)
		binary.LittleEndian.PutUint64(encrypted, uint64(c))
		result = append(result, encrypted...)
	}
	return result
}

func DecryptBuffer(enc []byte, priv *key.PrivateKey) []byte {
	var result []byte
	for i := 0; i < len(enc); i += 8 {
		g := binary.LittleEndian.Uint64(enc[i : i+8])
		c := Decrypt(int64(g), priv)
		result = append(result, byte(c))
	}
	return result
}

func StartRSA() {
	private, public, _ := CreateKeys()
	e := Encrypt(5, public)
	println("E=", e)
	m := Decrypt(e, private)
	println("M'=", m)
	src := []byte{'t', 'e', 's', 't', '2'}
	_ = ioutil.WriteFile("rsa-source.txt", src, os.ModePerm)
	enc := EncryptBuffer(src, public)
	_ = ioutil.WriteFile("rsa-encrypted.txt", enc, os.ModePerm)
	dec := DecryptBuffer(enc, private)
	_ = ioutil.WriteFile("rsa-decrypted.txt", dec, os.ModePerm)
}
