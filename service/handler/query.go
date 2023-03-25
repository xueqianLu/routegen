package handler

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xueqianLu/routegen/log"
	"github.com/xueqianLu/routegen/service/backend"
	"github.com/xueqianLu/routegen/service/param"
)

type RouteQuery struct {
	BaseController
}

func (q *RouteQuery) Route() {
	log.Debugf("got query")
	var query param.QueryRouteParam
	data := q.Ctx.Input.RequestBody
	if err := json.Unmarshal(data, &query); err != nil {
		logs.Error(err)
		q.ResponseInfo(500, "parse param failed", nil)
		return
	}
	result := backend.QueryRoute(query)
	q.ResponseInfo(200, nil, result)
}

func (q *RouteQuery) MergedRoute() {
	var query param.QueryRouteParam
	data := q.Ctx.Input.RequestBody
	if err := json.Unmarshal(data, &query); err != nil {
		logs.Error(err)
		q.ResponseInfo(500, "parse param failed", nil)
		return
	}
	result := backend.QueryRoute(query)
	// todo: add route merge filter
	q.ResponseInfo(200, nil, result)
}

func (q *RouteQuery) Version() {
	q.ResponseInfo(200, nil, "1.0.0")
}
