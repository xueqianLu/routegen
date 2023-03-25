package models

import (
	"github.com/zhihu/norm"
)

type Token struct {
	norm.VModel
	Name    string `norm:"name"`
	Address string `norm:"address"`
}

type Pair struct {
	norm.EModel
	DexName       string `norm:"dex"`
	PairAddress   string `norm:"pairaddress"`
	TrackedVolume string `norm:"tracked"`
	Fee           string `norm:"fee"`
	Token0        string `norm:"token0"`
	Token1        string `norm:"token1"`
}

var _ norm.IVertex = new(Token)
var _ norm.IEdge = new(Pair)

func (*Token) TagName() string {
	return "token"
}

func (t *Token) GetVid() interface{} {
	return t.Address
}

func (p *Pair) EdgeName() string {
	return "pair"
	//return fmt.Sprintf("%s", p.PairAddress)
}
