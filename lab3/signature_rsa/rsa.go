package signature_rsa

import (
	"encoding/json"
	"github.com/lol1pop/info_security_labs/basic"
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa/key"
	"io/ioutil"
	"math/big"
	"os"
)

type SignMess struct {
	S   *big.Int       `json:"s"`
	Key *key.PublicKey `json:"key"`
}

func CreateKeys() (*key.PrivateKey, *key.PublicKey) {
	privateKey := key.CreatePrivate()
	publicKey := key.CreatePublic(privateKey)
	return privateKey, publicKey
}

func CreateSignMessage(m []byte, privateKey *key.PrivateKey) (*big.Int, []byte) {
	h := basic.GetMessageHash(m)
	s := new(big.Int).Exp(h, privateKey.C, privateKey.N)
	return s, m
}

func CheckSignMessage(sign *big.Int, m []byte, publicKey *key.PublicKey) bool {
	h := basic.GetMessageHash(m)
	e := new(big.Int).Exp(sign, publicKey.D, publicKey.N)
	return e.Cmp(h) == 0
}

func StartRsaSign() {
	m := "test message"
	_ = ioutil.WriteFile("RsaSign.txt", []byte(m), os.ModePerm)
	fsrc, _ := ioutil.ReadFile("RsaSign.txt")
	println(string(fsrc))
	println("=== RSA ===")
	private, public := CreateKeys()
	println("N ", private.N.String())
	println("C ", private.C.String())
	println("D ", public.D.String())
	sign, _ := CreateSignMessage(fsrc, private)
	toJson := func(any interface{}) []byte {
		bytes, _ := json.Marshal(any)
		return bytes
	}
	_ = ioutil.WriteFile("RsaSign-private.txt", toJson(private), os.ModePerm)
	_ = ioutil.WriteFile("RsaSign-public.txt", toJson(public), os.ModePerm)
	_ = ioutil.WriteFile("RsaSign-sing.txt", toJson(SignMess{sign, public}), os.ModePerm)
	println("Check result: ", CheckSignMessage(sign, fsrc, public))
	CheckSignFile()
}

func CheckSignFile() {
	fsrc, _ := ioutil.ReadFile("RsaSign.txt")
	src, _ := ioutil.ReadFile("RsaSign-sing.txt")
	var data *SignMess
	_ = json.Unmarshal(src, &data)
	println("Check result: ", CheckSignMessage(data.S, fsrc, data.Key))
}
