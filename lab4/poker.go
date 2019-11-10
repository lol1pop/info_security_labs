package poker

import (
	"crypto/rand"
	"fmt"
	"math/big"
	rand_math "math/rand"
	"strconv"
	"time"
)

const Q_BIT = 16

var ONE = big.NewInt(1)
var TOW = big.NewInt(2)
var ZERO = big.NewInt(0)

var suit = []string{"♣", "♠", "♥", "♦"}
var high = []string{"J", "Q", "K", "A"}

type Player struct {
	P             *big.Int
	Ci            *big.Int
	Di            *big.Int
	onHand        []string
	onHandEncrypt []*big.Int
	onHandDecrypt []*big.Int
}

func InitParams() (_, _ *big.Int) {
	q, err := rand.Prime(rand.Reader, Q_BIT)
	if err != nil {
		panic(err)
	}
	p := new(big.Int).Add(new(big.Int).Mul(q, TOW), ONE)
	return p, q
}

func GeneratedKey(p *big.Int) (Ci, Di *big.Int) {
	p_1 := new(big.Int).Sub(p, ONE)
	for {
		for {
			Di, _ = rand.Int(rand.Reader, p_1)
			if new(big.Int).GCD(nil, nil, Di, p_1).Cmp(ONE) == 0 {
				break
			}
		}
		Ci = new(big.Int).ModInverse(Di, p_1)
		check := new(big.Int).Mod(new(big.Int).Mul(Di, Ci), p_1)
		if check.Cmp(ONE) == 0 {
			return
		}
	}

}

func ConnectPlayer(p *big.Int) Player {
	Ci, Di := GeneratedKey(p)
	return Player{P: p, Ci: Ci, Di: Di, onHand: []string{}}
}

func Shuffle(arr []*big.Int) []*big.Int {
	rand_math.Seed(time.Now().UnixNano())
	rand_math.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}

func (p *Player) EncryptCart(k *big.Int) *big.Int {
	return new(big.Int).Exp(k, p.Ci, p.P)
}

func (p *Player) DecryptCart(k *big.Int) *big.Int {
	return new(big.Int).Exp(k, p.Di, p.P)
}

func (p *Player) EncryptCartDeck(cart_deck []*big.Int) []*big.Int {
	var encrypt_arr []*big.Int
	for _, k := range cart_deck {
		encrypt := new(big.Int).Exp(k, p.Ci, p.P)
		encrypt_arr = append(encrypt_arr, encrypt)
	}
	return Shuffle(encrypt_arr)
}

func InitDeck(poker bool) (cart_deck []string) {
	var intiNum int
	if poker {
		intiNum = 2
	} else {
		intiNum = 6
	}

	for _, v := range suit {
		for i := intiNum; i <= 10; i++ {
			cart := strconv.FormatInt(int64(i), 10) + v
			cart_deck = append(cart_deck, cart)
		}
		for _, h := range high {
			cart := h + v
			cart_deck = append(cart_deck, cart)
		}
	}
	return
}
func AssociatedBigIntArr(len int) []*big.Int {
	var big_arr []*big.Int
	for i := 0; i < len; i++ {
		big_arr = append(big_arr, big.NewInt(int64(i)))
	}
	return big_arr
}

func StartPoker() {
	p, q := InitParams()
	println("========🃑\n P:", p.String(), "\n Q:", q.String(), "\n========🃑\n")
	fmt.Println(Shuffle([]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5)}))
}
