package poker

import (
	"crypto/rand"
	"fmt"
	"math/big"
	rand_math "math/rand"
	"strconv"
	"time"
)

const Q_BIT = 8

var ONE = big.NewInt(1)
var TOW = big.NewInt(2)
var ZERO = big.NewInt(0)

var suit = []string{"♣", "♠", "♥", "♦"}
var high = []string{"J", "Q", "K", "A"}

type Player struct {
	P             *big.Int
	Ci            *big.Int
	Di            *big.Int
	OnHand        []string
	OnHandEncrypt []*big.Int
	OnHandDecrypt []*big.Int
}

func InitParams() (_, _ *big.Int) {
	for {
		q, err := rand.Prime(rand.Reader, Q_BIT)
		if err != nil {
			panic(err)
		}
		p := new(big.Int).Add(new(big.Int).Mul(q, TOW), ONE)
		if p.ProbablyPrime(50) {
			return p, q
		}
	}
}

func GeneratedKey(p *big.Int) (Ci, Di *big.Int) {
	p_1 := new(big.Int).Sub(p, ONE)
	for {
		for {
			Di, _ = rand.Int(rand.Reader, p_1)
			check := new(big.Int).GCD(nil, nil, Di, p_1)
			if check.Cmp(ONE) == 0 {
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
	return Player{P: p, Ci: Ci, Di: Di}
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
	arr := Shuffle([]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5)})

	Alisa := ConnectPlayer(p)
	fmt.Println("Alisa:", Alisa)
	Bob := ConnectPlayer(p)
	fmt.Println("Bob:", Bob)
	Eva := ConnectPlayer(p)
	fmt.Println("Eva:", Eva)

	cart := big.NewInt(5)
	fmt.Println("cart: ", cart.String())
	cart = Alisa.EncryptCart(cart)
	fmt.Println("Alisa encrypt:", cart.String())
	cart = Bob.EncryptCart(cart)
	fmt.Println("Bob encrypt:", cart.String())
	cart = Eva.EncryptCart(cart)
	fmt.Println("Eva encrypt:", cart.String())
	cart = Alisa.DecryptCart(cart)
	fmt.Println("Alisa Decrypt:", cart.String())
	cart = Eva.DecryptCart(cart)
	fmt.Println("Eva Decrypt:", cart.String())
	cart = Bob.DecryptCart(cart)
	fmt.Println("Bob Decrypt:", cart.String())

	//deck := InitDeck(true)
	fmt.Println(arr)
	//ass_arr := AssociatedBigIntArr(len(deck))
	players := GeneratedPlayers(p, 3)
	encryptDeck := EncryptDeckAllPlayers(players, arr)
	fmt.Println(players)
	fmt.Println(encryptDeck)
	allotCartsPlayers(&players, &encryptDeck, 4)
	fmt.Println(players)
	fmt.Println(encryptDeck)

}

func GeneratedPlayers(p *big.Int, n int) []Player {
	var players []Player
	for i := 0; i < n; i++ {
		player := ConnectPlayer(p)
		players = append(players, player)
		fmt.Println("player N", i, ": ", player)
	}
	return players
}

func EncryptDeckAllPlayers(players []Player, encrypt []*big.Int) []*big.Int {
	for i, player := range players {
		encrypt = player.EncryptCartDeck(encrypt)
		fmt.Println("Encrypt player N", i, ": ", encrypt)
	}
	return encrypt
}

func allotCart(player *Player, deck *[]*big.Int) {
	num := len(*deck) - 1
	if num < 0 {
		fmt.Println()
		return
	}
	cart := (*deck)[num]
	player.OnHandEncrypt = append(player.OnHandEncrypt, cart)
	*deck = (*deck)[:num]
	fmt.Println(cart.String())
}

func allotCartsPlayers(players *[]Player, deck *[]*big.Int, n int) {
	for i := 0; i < n; i++ {
		for j, _ := range *players {
			fmt.Print("get catr player N", j, ": ")
			allotCart(&(*players)[j], deck)
		}
	}
}
