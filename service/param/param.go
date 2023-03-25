package param

import "github.com/xueqianLu/routegen/types"

type QueryRouteParam struct {
	Token0 string `json:"token0"`
	Token1 string `json:"token1"`
}

type QueryRouteResponse struct {
	Routes []*types.TokenRoute `json:"routes"`
}
