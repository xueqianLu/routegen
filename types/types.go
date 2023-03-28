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
	"encoding/json"
	"fmt"
)

type RoutePairInfo struct {
	Pair string `json:"pair"`
	Fee  string `json:"fee"`
	Dex  string `json:"dex"`
}

func TextAddress(addr string) string {
	return fmt.Sprintf("\"%s\"", addr)
}

type RouteStep struct {
	Pairs []RoutePairInfo `json:"pair"`
	Src   string          `json:"from"`
	Dst   string          `json:"to"`
}

type TokenRoute struct {
	Steps []RouteStep `json:"steps"`
}

func (r TokenRoute) String() string {
	d, _ := json.Marshal(r)
	return string(d)
}

type SortTokenRoutes []*TokenRoute

func (s SortTokenRoutes) Len() int           { return len(s) }
func (s SortTokenRoutes) Less(i, j int) bool { return len(s[i].Steps) < len(s[j].Steps) }
func (s SortTokenRoutes) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
