package web

import (
	"mmschainnewaccount/config"
	routers "mmschainnewaccount/web/routers"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/astaxie/beego"
)


func Start() {

	beego.BConfig.WebConfig.Session.SessionOn = false
	beego.BConfig.Listen.HTTPPort = int(config.WebPort)

	beego.BConfig.Listen.HTTPSAddr = config.WebAddr

	beego.BConfig.WebConfig.Session.SessionName = "mmschainnewaccount"

	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.TemplateLeft = "<%"
	beego.BConfig.WebConfig.TemplateRight = "%>"










	beego.BConfig.WebConfig.ViewsPath = config.Web_path_views
	beego.SetStaticPath("/static", config.Web_path_static)
	routers.Start()


}


func openLocalWeb() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start http:
	case "darwin":
		cmd = exec.Command("open", "http:
	case "linux":
		cmd = exec.Command("xdg-open", "http:

	}
	err := cmd.Start()
	if err != nil {

	}
}
