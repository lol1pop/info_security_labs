package course

import (
	"crypto/rand"
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa"
	"math/big"
	rand_math "math/rand"
	"time"
)

// ============== func Alisa ================

func Shuffle(arr []int) []int {
	rand_math.Seed(time.Now().UnixNano())
	rand_math.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}

func PlanColorBit(c *big.Int, color int) *big.Int {
	if color == 10 {
		c = new(big.Int).Or(c, new(big.Int).Lsh(big.NewInt(1), 1))
		return new(big.Int).And(c, new(big.Int).Not(new(big.Int).Lsh(big.NewInt(1), 0)))
	}
	if color == 0 {
		c = new(big.Int).And(c, new(big.Int).Not(new(big.Int).Lsh(big.NewInt(1), 0)))
		return new(big.Int).And(c, new(big.Int).Not(new(big.Int).Lsh(big.NewInt(1), 1)))
	}

	c = new(big.Int).And(c, new(big.Int).Not(new(big.Int).Lsh(big.NewInt(1), 1)))
	return new(big.Int).Or(c, new(big.Int).Lsh(big.NewInt(1), 0))
}

func generatedR(graph *[][]Node, colorArr []int) {
	for i := 0; i < len(*graph); i++ {
		for j := 0; j < len((*graph)[i]); j++ {
			color := colorArr[(*graph)[i][j].Color]
			c, _ := rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(16), nil))
			r := PlanColorBit(c, color)
			(*graph)[i][j].R = r
		}
	}
}

func encryptRsaNode(graph *[][]Node) {
	for i := 0; i < len(*graph); i++ {
		for j := 0; j < len((*graph)[i]); j++ {
			(*graph)[i][j].PrivateKey, (*graph)[i][j].PublicKey = signature_rsa.CreateKeys()
		}
	}
}

func CreatedZu(graph *[][]Node) [][]EncryptNode {
	var encryptArrayNodes [][]EncryptNode
	for i := 0; i < len(*graph); i++ {
		var encryptNodes []EncryptNode
		for j := 0; j < len((*graph)[i]); j++ {
			encryptNodes = append(encryptNodes, EncryptNode{
				N:    (*graph)[i][j].PublicKey.N,
				D:    (*graph)[i][j].PublicKey.D,
				Z:    new(big.Int).Exp((*graph)[i][j].R, (*graph)[i][j].PublicKey.D, (*graph)[i][j].PublicKey.N),
				Edge: (*graph)[i][j].Edge,
			})
		}
		encryptArrayNodes = append(encryptArrayNodes, encryptNodes)
	}
	return encryptArrayNodes
}

func getDecryptKey(graph [][]Node, edge int) (*big.Int, *big.Int) {
	nodes := graph[edge]
	return nodes[0].PrivateKey.C, nodes[1].PrivateKey.C
}

//========== func BOB ===========

func getRandomEdge(arr []int) int {
	return rand_math.Intn(len(arr))
	//arr[:len(arr) - 1]
}

func getColorFromR(r *big.Int) int {
	checkTen := new(big.Int).And(r, new(big.Int).Lsh(big.NewInt(1), 1)).Uint64()
	checkOne := new(big.Int).And(r, new(big.Int).Lsh(big.NewInt(1), 0)).Uint64()
	if checkTen == 2 {
		return 10
	} else {
		if checkOne == 1 {
			return 1
		} else {
			return 0
		}
	}
}

func checkColors(encryptNodes [][]EncryptNode, edge int, c1 *big.Int, c2 *big.Int) {
	u := encryptNodes[edge]
	u1 := u[0]
	IZu1 := new(big.Int).Exp(u1.Z, c1, u1.N)
	u2 := u[0]
	IZu2 := new(big.Int).Exp(u2.Z, c2, u2.N)
	color1 := getColorFromR(IZu1)
	color2 := getColorFromR(IZu2)
	print(color1)
	print(" ")
	println(color2)
	println(color1 != color2)
}

func Start() {
	graph := [][]Node{
		{Node{
			Edge:       0,
			Color:      0,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}, Node{
			Edge:       0,
			Color:      2,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}},
		{Node{
			Edge:       1,
			Color:      2,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}, Node{
			Edge:       1,
			Color:      0,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}},
		{Node{
			Edge:       2,
			Color:      0,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}, Node{
			Edge:       2,
			Color:      1,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}},
		{Node{
			Edge:       3,
			Color:      2,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}, Node{
			Edge:       3,
			Color:      1,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}},
		{Node{
			Edge:       4,
			Color:      2,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}, Node{
			Edge:       4,
			Color:      0,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}},
		{Node{
			Edge:       5,
			Color:      1,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}, Node{
			Edge:       5,
			Color:      0,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}},
		{Node{
			Edge:       6,
			Color:      0,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}, Node{
			Edge:       6,
			Color:      2,
			R:          nil,
			PrivateKey: nil,
			PublicKey:  nil,
		}},
	}
	colorArr := []int{0, 1, 10} //R,B,Y
	generatedR(&graph, Shuffle(colorArr))
	encryptRsaNode(&graph)
	encryptNodes := CreatedZu(&graph)
	edge := getRandomEdge([]int{0, 1, 2, 3, 4, 5, 6})
	c1, c2 := getDecryptKey(graph, edge)
	checkColors(encryptNodes, edge, c1, c2)
}
