package diffie_hellman

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

func PowByModule(a, x, p uint64) uint64 {
	var r uint64 = 1
	for i := x; i != 0; i = i >> 1{
		if i&1 == 1{
			r = r * a % p
		}
		a = a * a % p
	}
	return r
}

type Uravnenie struct {
	gcd uint64
	x   uint64
	y   uint64
}

func Gcd(a, b uint64) uint64 {
	var r uint64
	for ; b != 0; {
		r = b
		b = a % b
		a = r
	}
	return r
}

func CommonGcd(a, b uint64) (_, _, _ uint64) {
	u, v := Uravnenie{a, 1, 0}, Uravnenie{b, 0, 1}
	for v.gcd != 0 {
		q := u.gcd / v.gcd
		t := Uravnenie{u.gcd % v.gcd, u.x - q*v.x, u.y - q*v.y}
		u, v = v, t
	}
	return u.gcd, u.x, u.y
}

func RangeRandom(left, right uint64) uint64 {
	min := int64(left)
	max := int64(right) - min
	return uint64(rand.Int63n(max))
}

func TestFerma(p uint64, k int) bool {
	if p == 2 { return true }
	if p&1 != 1 { return false }
	for i:= 0; i < k; i++ {
		a := uint64(rand.Int63n(int64(p) - 3) + 3)
		if Gcd(a, p) != 1 || PowByModule(a, p -1 , p) != 1 {
			return false
		}
	}
	return true
}

func DiffieHellmanGeneratedDigit() (_,_ uint64) {
	for {
		q := RangeRandom(50,1000000000)

		fmt.Print("|:")
		fmt.Println(q)
		//q := rand.Uint64()
		if TestFerma(q, 15) {
			fmt.Print("  -->: ")
			fmt.Println(q)
			p := 2 * q + 1
			if TestFerma(p, 15){
				return q, p
			}
		}
	}
}

func DiffieHellmanGeneratedFreeData() (_,_ uint64) {
	q, p := DiffieHellmanGeneratedDigit()
	for {
		max := int64(p) - 1
		g := uint64(rand.Int63n(max))
		if 1 < g && g < (p - 1) {
			if PowByModule(g, q, p) != 1 {
				return g, p
			}
		}
	}
}

func DiffieHellmanCreatedKey(a, p uint64) (public, private uint64) {
	key := RangeRandom(50,1000000000)
	//key := rand.Uint64()
	publicKey := PowByModule(a, key, p)
	return publicKey, key
}

func DiffieHellmanPing() (requestData Uravnenie, private uint64) {
	a, p := DiffieHellmanGeneratedFreeData()
	publicKey, key := DiffieHellmanCreatedKey(a, p)
	return Uravnenie{x: a , y: p, gcd: publicKey}, key
}

func DiffieHellmanPong(keys Uravnenie) (public, cipher uint64) {
	publicKey, key := DiffieHellmanCreatedKey(keys.x, keys.y)
	cipher = DiffieHellmanCipher(keys.gcd, key, keys.y)
	return publicKey, cipher
}

func DiffieHellmanCipher(publicKey, privateKey, p uint64) uint64 {
	return PowByModule(publicKey, privateKey, p)
}

func Steps(a, p, y uint64) (uint64, error) {
	sqrt := uint64(math.Sqrt(float64(p)))
	m, k := sqrt, sqrt+1
	for m*k <= p {
		k++
	}
	vals := make(map[uint64]uint64)
	counter := 0
	for j := uint64(0); j < m; j++ {
		t := (PowByModule(a, j, p) * y) % p
		vals[t] = j
		counter++
	}
	for i := uint64(1); i <= k; i++ {
		counter++
		t := PowByModule(a, i*m, p)
		if j, ok := vals[t]; ok {
			return i*m - j, nil
		}
	}
	println("Steps: ", counter)
	return 0, fmt.Errorf("Invalid in data")
}

func StartDiffiehellman() {
	//	rand.Seed(time.Now().UnixNano())
	//fmt.Println("2^10 mod 5: ", PowByModule(2,10, 5))
	//fmt.Println("3^21 mod 11: ", PowByModule(3,21, 11))
	//fmt.Println("7^31 mod 17: ", PowByModule(7,31, 17))
	//fmt.Println("5^17 mod 11: ", PowByModule(5,17, 11))
	//fmt.Println("3^15 mod 10: ", PowByModule(3,15, 10))
	//fmt.Println("7^12 mod 9: ", PowByModule(2,9, 25))
	//fmt.Println("28 Gcd 2: ", Gcd(28,8))
	//fmt.Println("28 Gcd 2: ", Gcd(28,2))
	//fmt.Println("28 Gcd 19: ", Gcd(28,19))
	//fmt.Println("19 Gcd 9: ", Gcd(19,9))
	//fmt.Println("9 Gcd 1: ", Gcd(9,1))

	//	rand.Seed(time.Now().UnixNano())
	fmt.Println("===============Diffie Hellman=================")
	fmt.Println("ALISA  --ping--> BOB : ")
	requestData, keyAlisa := DiffieHellmanPing()
	fmt.Println(requestData)
	fmt.Println("keyAlisa : " + strconv.FormatUint(keyAlisa,10))
	fmt.Println("BOB ---get request--> ALISA")
	publicBOB, cipherBOB := DiffieHellmanPong(requestData)
	fmt.Println("BOB  --pong--> ALISA : " + strconv.FormatUint(publicBOB, 10))
	fmt.Println("publicBOB : " + strconv.FormatUint(publicBOB,10))
	fmt.Println("ALISA ---get pong--> BOB ")
	cipherALISA := DiffieHellmanCipher(publicBOB, keyAlisa, requestData.y)
	fmt.Println("cipher ALISA :" + strconv.FormatUint(cipherALISA, 10) + "  |||  cipher BOB :" + strconv.FormatUint(cipherBOB, 10))
	println(cipherALISA == cipherBOB)

	//fmt.Println("===============Baby step / gigant step=================")
	//d := PowByModule(2,123456,1000000007)
	//fmt.Println(d)
	//x, err := Steps(2, 1000000007, d)
	//if err == nil {
	//	println(x)
	//	result := PowByModule(2, x, 1000000007)
	//	println(3 == result)
	//} else {
	//	println(err.Error())
	//}
}