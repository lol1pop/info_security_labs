package course

import (
	"crypto/rand"
	"encoding/json"
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa"
	"io/ioutil"
	"math/big"
	rand_math "math/rand"
	"os"
	"strconv"
	"time"
)

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

type Sender struct {
	graph  [][]Node
	colors []int
}

func (s *Sender) init() *Sender {
	gr, _ := ioutil.ReadFile("graph.json")
	fromJson := func(any []byte) [][]Node {
		var graph [][]Node
		if err := json.Unmarshal(any, &graph); err != nil {
			panic(err)
		}
		return graph
	}
	s.graph = fromJson(gr)
	s.colors = []int{0, 1, 10} //R,B,Y
	return s
}

func (s *Sender) repainting() {
	s.colors = Shuffle(s.colors)
}

func (s *Sender) generatedR() {
	s.repainting()
	for i := 0; i < len(s.graph); i++ {
		for j := 0; j < len(s.graph[i]); j++ {
			color := s.colors[s.graph[i][j].Color]
			c, _ := rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(16), nil))
			r := PlanColorBit(c, color)
			s.graph[i][j].R = r
		}
	}
}

func (s *Sender) encryptRsaNodes() {
	for i := 0; i < len(s.graph); i++ {
		for j := 0; j < len(s.graph[i]); j++ {
			s.graph[i][j].PrivateKey, s.graph[i][j].PublicKey = signature_rsa.CreateKeys()
		}
	}
}

func (s Sender) createdZu() [][]EncryptNode {
	var encryptArrayNodes [][]EncryptNode
	for i := 0; i < len(s.graph); i++ {
		var encryptNodes []EncryptNode
		for j := 0; j < len(s.graph[i]); j++ {
			encryptNodes = append(encryptNodes, EncryptNode{
				N:    s.graph[i][j].PublicKey.N,
				D:    s.graph[i][j].PublicKey.D,
				Z:    new(big.Int).Exp(s.graph[i][j].R, s.graph[i][j].PublicKey.D, s.graph[i][j].PublicKey.N),
				Edge: s.graph[i][j].Edge,
			})
		}
		encryptArrayNodes = append(encryptArrayNodes, encryptNodes)
	}
	return encryptArrayNodes
}

func (s Sender) getDecryptKey(edge int) (*big.Int, *big.Int) {
	nodes := s.graph[edge]
	return nodes[0].PrivateKey.C, nodes[1].PrivateKey.C
}

type Receiver struct {
	edgesArr     []int
	encryptNodes [][]EncryptNode
}

func (r *Receiver) init(size int) *Receiver {
	for i := 0; i < size; i++ {
		r.edgesArr = append(r.edgesArr, i)
	}
	return r
}

func (r *Receiver) getRandEdge() int {
	edges := Shuffle(r.edgesArr)
	edge := edges[0]
	r.edgesArr = edges[1:]
	return edge
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

func (r *Receiver) checkColors(edge int, c1 *big.Int, c2 *big.Int) {
	u := r.encryptNodes[edge]
	u1 := u[0]
	IZu1 := new(big.Int).Exp(u1.Z, c1, u1.N)
	u2 := u[1]
	IZu2 := new(big.Int).Exp(u2.Z, c2, u2.N)
	color1 := getColorFromR(IZu1)
	color2 := getColorFromR(IZu2)
	print(color1)
	print(" ")
	println(color2)
	println(color1 != color2)
}

func (r *Receiver) updateEncryptNodes(encryptNodes [][]EncryptNode) {
	r.encryptNodes = encryptNodes
}

func Start() {
	sender := new(Sender).init()
	receiver := new(Receiver).init(len(sender.graph))

	for len(receiver.edgesArr) > 0 {
		sender.generatedR()
		sender.encryptRsaNodes()
		receiver.updateEncryptNodes(sender.createdZu())

		edge := receiver.getRandEdge()
		print("[" + strconv.Itoa(edge) + "]: ")
		c1, c2 := sender.getDecryptKey(edge)
		receiver.checkColors(edge, c1, c2)
	}
}

func CreatedGraph() {
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
	toJson := func(any interface{}) []byte {
		bytes, _ := json.Marshal(any)
		return bytes
	}
	_ = ioutil.WriteFile("graph.json", toJson(graph), os.ModePerm)
}
