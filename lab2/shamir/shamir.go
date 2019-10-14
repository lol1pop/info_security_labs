package shamir

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/lol1pop/info_security_labs/basic"
	"io"
	"io/ioutil"
	"os"
)

func GeneratedPrimeNumber(left, right uint64) uint64 {
	for {
		p := basic.RangeRandom(left, right)
		//fmt.Print("|:")
		//fmt.Println(p)
		if basic.TestFerma(p, 15) {
			//fmt.Print("  -->: ")
			//fmt.Println(p)
			return p
		}
	}
}

func inversionDigit(F uint64, Di uint64) uint64 {
	_, _, invert := basic.CommonGcd(int64(F), int64(Di))
	if invert < 0 {
		fi := int64(F)
		invert = fi + invert
	}
	return uint64(invert)
}

type User struct {
	Name        string
	EncryptKey  uint64
	DecryptKey  uint64
	PrimeNumber uint64
}

func (k *User) Save(wr io.Writer) error {
	bytes, err := json.Marshal(*k)
	if err == nil {
		_, err = wr.Write(bytes)
	}
	return err
}

func LoadKeyFromReader(filename string) (User, error) {
	var marshalled []byte
	var err error
	if marshalled, err = ioutil.ReadFile(filename); err != nil {
		return User{}, err
	}
	var key User
	err = json.Unmarshal(marshalled, &key)
	return key, err
}

func initShamir() uint64 {
	return GeneratedPrimeNumber(1000000, 10000000)
}

func InitUserShamir(name string, primeNumber uint64) User {
	save := func(key User, filename string) error {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		return key.Save(file)
	}
	encryptKey, decryptKey := shamirGeneratedDigit(primeNumber)
	user := User{
		Name:        name,
		PrimeNumber: primeNumber,
		EncryptKey:  encryptKey,
		DecryptKey:  decryptKey,
	}
	_ = save(user, "keys.shamir")
	return user
}

/*
* Число D_i_  должно быть взаимно простое с числом р-1
 */
func generatedDi(p uint64) uint64 {
	for {
		genDigit := basic.RangeRandom(50, p-1)
		if basic.Gcd(genDigit, p-1) == 1 {
			return genDigit
		}
	}
}

func shamirGeneratedDigit(p uint64) (c, d uint64) {
	for {
		Di := generatedDi(p)
		Ci := inversionDigit(p-1, Di)
		chec := (Di * Ci) % (p - 1)
		if chec == 1 {
			return Ci, Di
		}
	}
}

func EncryptMessage(message []byte, user User) []byte {
	var encryptMessage []byte
	encryptSymbol := make([]byte, 8)
	for _, m := range message {
		x := basic.PowByModule(uint64(m), user.EncryptKey, user.PrimeNumber)
		binary.LittleEndian.PutUint64(encryptSymbol, x)
		encryptMessage = append(encryptMessage, encryptSymbol...)
	}
	return encryptMessage
}

func DecryptMessage(message []byte, user User) []byte {
	var decryptMessage []byte
	for i := 0; i < len(message); i += 8 {
		m := binary.LittleEndian.Uint64(message[i : i+8])
		x := basic.PowByModule(m, user.DecryptKey, user.PrimeNumber)
		decryptMessage = append(decryptMessage, byte(x))
	}
	return decryptMessage
}

func StartShamir() {
	p := initShamir()
	alisa := InitUserShamir("Alisa", p)
	falisa, _ := LoadKeyFromReader("keys.shamir")
	eMgs := EncryptMessage([]byte{'t'}, falisa)
	println("E=", string(eMgs))
	dMgs := DecryptMessage(eMgs, falisa)
	println("D=", string(dMgs))
	fmt.Println(alisa)
	src := []byte{'t', 'e', 's', 't', '1'}
	_ = ioutil.WriteFile("source-shamir.txt", src, os.ModePerm)
	fsrc, _ := ioutil.ReadFile("source-shamir.txt")
	enc := EncryptMessage(fsrc, alisa)
	_ = ioutil.WriteFile("encrypted-shamir.txt", enc, os.ModePerm)
	dec := DecryptMessage(enc, alisa)
	_ = ioutil.WriteFile("decrypted-shamir.txt", dec, os.ModePerm)

}
