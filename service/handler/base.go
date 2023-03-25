package handler

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (d *BaseController) ResponseInfo(code int, errMsg interface{}, result interface{}) {
	//codestr := strconv.Itoa(code)
	//d.Data["json"] = map[string]interface{}{"error": codestr, "err_msg": errMsg, "data": result}
	switch code {
	case 200:
		d.Data["json"] = map[string]interface{}{"code": 200, "err_msg": errMsg, "data": result}
	default:
		d.Data["json"] = map[string]interface{}{"code": 500, "err_msg": errMsg}
	}
	d.ServeJSON()
}
