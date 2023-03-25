package service

import (
	"github.com/astaxie/beego"
	_ "github.com/xueqianLu/routegen/service/router"
)

func StartServier() {
	//conf := config.GetConfig()
	//router.InitRoute()
	beego.Run()
	//beego.Run(conf.ServerAddr)
}
