package backend

import (
	"errors"
	"github.com/xueqianLu/routegen/config"
	"github.com/xueqianLu/routegen/database"
	"github.com/xueqianLu/routegen/service/param"
	"github.com/zhihu/norm"
)

var (
	b *Backend
)

type Backend struct {
	db *norm.DB
}

func SetupBackend() error {
	b = new(Backend)
	db := database.NewDb(config.GetConfig())
	if db == nil {
		return errors.New("create db failed")
	}
	b.db = db
	return nil
}

func QueryRoute(query param.QueryRouteParam) *param.QueryRouteResponse {
	paths := database.QueryRoute(b.db, query.Token0, query.Token1)
	result := new(param.QueryRouteResponse)
	result.Routes = paths
	return result
}
