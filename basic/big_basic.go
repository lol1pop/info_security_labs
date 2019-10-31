package basic

import "math/big"

func BigFastPowByModule(a, x, p *big.Int) *big.Int {
	r := big.NewInt(1)
	h := new(big.Int).SetBytes(a.Bytes())
	for i := new(big.Int).SetBytes(x.Bytes()); i.Cmp(big.NewInt(0)) != 0; i = i.Rsh(i, 1) {
		//println("(r,i):" + r.Text(10) + "  " + i.Text(10))
		if new(big.Int).And(i, big.NewInt(1)).Cmp(big.NewInt(1)) == 0 {
			r.Mul(r, h).Mod(r, p)
			//println("(r,if):" + r.Text(10) + "   t" )
		}
		h.Mul(h, h).Mod(h, p)
		//println( "(a):"+ a.Text(10))
	}
	return r
}
