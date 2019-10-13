package basic

import (
	"math/rand"
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
	gcd int64
	x   int64
	y   int64
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

func CommonGcd(a, b int64) (_, _, _ int64) {
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
