package router

import (
	"github.com/astaxie/beego"
	"github.com/xueqianLu/routegen/log"
	"github.com/xueqianLu/routegen/service/handler"
)

func init() {
	//func InitRoute() {
	log.Info("init router")
	beego.Router("/defiroute/api/v1/route", &handler.RouteQuery{}, "post:Route")
	beego.Router("/defiroute/api/v1/mergedroute", &handler.RouteQuery{}, "post:MergedRoute")
	beego.Router("/defiroute/api/v1/version", &handler.RouteQuery{}, "get:Version")
}
