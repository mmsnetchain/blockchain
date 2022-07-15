package anonymousnet

import (
	"bytes"
	"fmt"
	"github.com/prestonTao/libp2parea/nodeStore"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/name"

	"mmschainnewaccount/proxyhttp"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {

	this.TplName = "mmschainnewaccount/index.tpl"
}

func (this *MainController) Agent() {
	url := this.Ctx.Input.Param(":splat")

	fmt.Println("", url)

	temp := strings.SplitN(url, "/", 2)
	url = "/" + url

	var id *nodeStore.TempId
	port := uint16(0)

	temp = strings.SplitN(temp[0], ":", 2)

	height := mining.GetHighestBlock()

	targetMhId := name.FindNameToNetRandOne(temp[0], height)
	if targetMhId == nil {
		fmt.Println("")
		fmt.Println("id")
		id := nodeStore.AddressFromB58String(temp[0])
		targetMhId = &id
	}

	fmt.Println("id", targetMhId.B58String())

	if targetMhId.B58String() == "" {
		return
	}

	id = nodeStore.NewTempId(targetMhId, targetMhId)

	if len(temp) > 1 {
		var err error
		p, err := strconv.Atoi(temp[1])
		if err != nil {

			return
		}
		port = uint16(p)
	}

	bodyData := make([]byte, 0)
	this.Ctx.Request.Body.Read(bodyData)

	request := proxyhttp.HttpRequest{
		Port:    port,
		Method:  this.Ctx.Request.Method,
		Header:  this.Ctx.Request.Header,
		Body:    bodyData,
		Url:     url,
		Cookies: this.Ctx.Request.Cookies(),
	}

	bs := proxyhttp.SendHttpRequest(id, request)

	if bs == nil {
		fmt.Println("")
		return
	}

	response := new(proxyhttp.HttpResponse)

	decoder := json.NewDecoder(bytes.NewBuffer(*bs))
	decoder.UseNumber()
	err := decoder.Decode(&response)
	if err != nil {
		fmt.Println(err)
		return
	}

	this.Ctx.ResponseWriter.WriteHeader(response.StatusCode)

	for key, value := range response.Header {
		this.Ctx.ResponseWriter.Header().Set(key, value[0])
		for i := 1; i < len(value); i++ {
			this.Ctx.ResponseWriter.Header().Add(key, value[i])
		}
	}

	for _, one := range response.Cookies {
		http.SetCookie(this.Ctx.ResponseWriter, one)
	}

	this.Ctx.ResponseWriter.Write(response.Body)

	return

}

func (this *MainController) AgentToo() {
	fmt.Println("")
}
