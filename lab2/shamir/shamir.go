package shamir

import (
	"github.com/lol1pop/info_security_labs/basic"
	"encoding/binary"
	"fmt"
)

func GeneratedPrimeNumber(left, right  uint64) uint64 {
	for {
		p := basic.RangeRandom(left, right)
		fmt.Print("|:")
		fmt.Println(p)
		if basic.TestFerma(p, 15) {
			fmt.Print("  -->: ")
			fmt.Println(p)
			return p
		}
	}
}

func inversionDigit(F uint64,Di uint64 ) uint64 {
	_,_, invert := basic.CommonGcd(int64(F),int64(Di))
	if invert < 0 {
		fi := int64(F)
		invert = fi + invert
	}
	return uint64(invert)
}

type User struct {
	name string
	encryptKey uint64
	decryptKey uint64
	primeNumber uint64
}

func initShamir() uint64 {
	return GeneratedPrimeNumber(1000000, 10000000)
}

func InitUserShamir(name  string, primeNumber uint64) User {
	encryptKey, decryptKey := shamirGeneratedDigit(primeNumber)
	return User{
		name: name,
		primeNumber: primeNumber,
		encryptKey: encryptKey,
		decryptKey: decryptKey,
	}
}
/*
* Число D_i_  должно быть взаимно простое с числом р-1
*/
func generatedDi(p uint64) uint64 {
	for {
		genDigit := GeneratedPrimeNumber(50, p - 1)
		if basic.Gcd(genDigit, p - 1) == 1 {
			return genDigit
		}
	}
}

func shamirGeneratedDigit(p uint64) (c,d uint64) {
	for {
		Di := generatedDi(p)
		Ci := inversionDigit(p - 1, Di)
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
		x := basic.PowByModule(uint64(m), user.encryptKey, user.primeNumber)
		binary.LittleEndian.PutUint64(encryptSymbol, x)
		encryptMessage = append(encryptMessage, encryptSymbol...)
	}
	return encryptMessage
}

func DecryptMessage(message []byte,user User) []byte {
	var decryptMessage []byte
	for i := 0; i < len(message); i += 8 {
		m := binary.LittleEndian.Uint64(message[i : i+8])
		x := basic.PowByModule(m, user.decryptKey, user.primeNumber)
		decryptMessage = append(decryptMessage, byte(x))
	}
	return decryptMessage
}

func StartShamir()  {
	p :=  initShamir()
	alisa := InitUserShamir("Alisa", p)
    eMgs := EncryptMessage([]byte{'t'}, alisa)
    println("E=", string(eMgs))
    dMgs := DecryptMessage(eMgs,alisa)
	println("D=", string(dMgs))
	fmt.Println(alisa)
}