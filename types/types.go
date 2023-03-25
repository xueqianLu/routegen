/*
Copyright © 2023 xueqianLu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Token struct {
	Name     string         `json:"name"`
	Contract common.Address `json:"contract"`
}

type Pair struct {
	TokenIn       Token    `json:"token_in"`
	TokenOut      Token    `json:"token_out"`
	TrackedVolume *big.Int `json:"tracked_volume"`
}

func (p Pair) String() string {
	return fmt.Sprintf("%s_%s", p.TokenIn.Name, p.TokenOut.Name)
}

type Dex struct {
	Name  string  `json:"name"`
	Pairs []*Pair `json:"pairs"`
}

type RouteOne struct {
	RouteDex  *Dex  `json:"route_dex"`
	RoutePair *Pair `json:"route_pair"`
}

type Route struct {
	TokenIn   Token       `json:"token_in"`
	TokenOut  Token       `json:"token_out"`
	RoutePath []*RouteOne `json:"route_path"`
}
