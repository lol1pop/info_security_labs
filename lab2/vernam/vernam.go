package vernam

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func cipher(msg, key []byte) ([]byte, error) {
	if len(msg) != len(key) || len(msg) == 0 {
		return nil, fmt.Errorf("Invalid message and key size")
	}
	result := make([]byte, len(msg))
	for i := range msg {
		result[i] = msg[i] ^ key[i]
	}
	return result, nil
}

//Return: key, encrypted message
func CreateKeyAndEncrypt(msg []byte) ([]byte, []byte, error) {
	key := CreateKey(uint64(len(msg)))
	enc, err := cipher(msg, key)
	return key, enc, err
}

func Encrypt(msg, key []byte) ([]byte, error) {
	return cipher(msg, key)
}

func Decrypt(enc, key []byte) ([]byte, error) {
	return cipher(enc, key)
}

func VernamExample() {
	src := []byte{'t', 'e', 's', 't', '2'}
	key, enc, err := CreateKeyAndEncrypt(src)
	if err != nil {
		log.Printf("Encrypt error: %s", err.Error())
		return
	}
	println(string(enc))
	dec, err := Decrypt(enc, key)
	if err != nil {
		log.Printf("Decrypt error: %s", err.Error())
		return
	}
	println(string(dec))
	_ = ioutil.WriteFile("vernam-source.txt", src, os.ModePerm)
	_ = ioutil.WriteFile("vernam-key.txt", key, os.ModePerm)
	_ = ioutil.WriteFile("vernam-encrypt.txt", enc, os.ModePerm)
	_ = ioutil.WriteFile("vernam-decrypt.txt", dec, os.ModePerm)
}
