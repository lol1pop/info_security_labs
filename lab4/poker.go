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
var cartsdeck []string

type Player struct {
	Number        int
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

func ConnectPlayer(p *big.Int, num int) Player {
	Ci, Di := GeneratedKey(p)
	return Player{Number: num, P: p, Ci: Ci, Di: Di}
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

func (p *Player) DecryptCartDeck(cart_deck []*big.Int) []*big.Int {
	var decrypt_arr []*big.Int
	for _, k := range cart_deck {
		decrypt := new(big.Int).Exp(k, p.Di, p.P)
		decrypt_arr = append(decrypt_arr, decrypt)
	}
	return decrypt_arr
}

func (p *Player) DecryptHandDeck(players []Player) {
	for _, player := range players {
		if player.Number == p.Number {
			continue
		}
		fmt.Print("player N", p.Number, " decrypt carts by N", player.Number, ":")
		p.OnHandEncrypt = player.DecryptCartDeck(p.OnHandEncrypt)
		fmt.Println(p.OnHandEncrypt)
	}
	p.OnHandDecrypt = p.DecryptCartDeck(p.OnHandEncrypt)
	p.OnHandEncrypt = p.OnHandEncrypt[:0]
	fmt.Println("self decrypt:", p.OnHandDecrypt)
}

func (p *Player) AssociatedNumCartShirt() {
	for _, c := range p.OnHandDecrypt {
		cart := cartsdeck[c.Int64()]
		p.OnHand = append(p.OnHand, cart)
	}
	fmt.Println("player N", p.Number, " have Carts on hand:", p.OnHand)
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

func GeneratedPlayers(p *big.Int, n int) []Player {
	var players []Player
	for i := 0; i < n; i++ {
		player := ConnectPlayer(p, i)
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
			if len(*deck)-1 < 0 {
				fmt.Println()
				return
			}
			fmt.Print("get catr player N", j, ": ")
			allotCart(&(*players)[j], deck)
		}
	}
}

func PlayersDecryptHandDeck(players *[]Player) {
	for j, _ := range *players {
		(*players)[j].DecryptHandDeck(*players)
		(*players)[j].AssociatedNumCartShirt()
	}
}

func DecryptTableCart(players []Player, cart *big.Int) *big.Int {
	for _, player := range players {
		cart = player.DecryptCart(cart)
	}
	return cart
}

func DecryptTableCarts(players []Player, deck *[]*big.Int, n int) []*big.Int {
	var carts []*big.Int
	for i := 0; i < n; i++ {
		num := len(*deck) - 1 - i
		carts = append(carts, DecryptTableCart(players, (*deck)[i]))
		*deck = (*deck)[:num]
	}
	return carts
}

func ShowCardsTable(players []Player, deck *[]*big.Int, n int) {
	println("Carts on table : ")
	var show []string
	carts := DecryptTableCarts(players, deck, n)
	for _, c := range carts {
		cart := cartsdeck[c.Int64()]
		show = append(show, cart)
		print(cart, " ")
	}
	println()
}

func StartPoker() {
	nP := 3
	nC := 4
	p, q := InitParams()
	fmt.Println("========🃑\n P:", p.String(), "\n Q:", q.String(), "\n========🃑\n")
	cartsdeck = InitDeck(true)
	ass_arr := AssociatedBigIntArr(len(cartsdeck))
	players := GeneratedPlayers(p, nP)
	encryptDeck := EncryptDeckAllPlayers(players, ass_arr)
	fmt.Println(players)
	fmt.Println(encryptDeck)
	allotCartsPlayers(&players, &encryptDeck, nC)
	fmt.Println(players)
	fmt.Println(encryptDeck)
	PlayersDecryptHandDeck(&players)
	fmt.Println(players)

	println("Desk: ", nP, " players, rule allot by", nC, " carts")
	println("Deсk: ", len(cartsdeck), " carts")
	for _, player := range players {
		println("Player#", player.Number, " on hand:")
		for _, c := range player.OnHand {
			println("  -->" + c)
		}
	}
	ShowCardsTable(players, &encryptDeck, 5)
}
