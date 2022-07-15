package rpc

import (
	"io/ioutil"
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc/model"
	"net"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
)

var (
	Allowip = ""
)

type Handler struct {
	w    http.ResponseWriter
	r    *http.Request
	body []byte
}

func (h *Handler) init(w http.ResponseWriter, r *http.Request) *Handler {
	h.w = w
	h.r = r
	return h
}
func (h *Handler) Out(data []byte) {
	h.w.Header().Add("Content-Type", "application/json")
	datas := append(append([]byte(`{"jsonrpc":"2.0","code":`), append([]byte(strconv.Itoa(model.Success)), byte(','))...), data[1:]...)
	h.w.Write(datas)
	return
}
func (h *Handler) Err(code, data string) {

	h.w.Header().Add("Content-Type", "application/json")
	h.w.Write([]byte(`{"jsonrpc":"2.0","code":` + code + `,"message":"` + data + `"}`))
	return
}
func (h *Handler) Validate() (msg string, ok bool) {

	if Allowip != "" && h.RemoteIp() != Allowip {
		msg = "deny ip"
		ok = true
	}

	if h.r.Header.Get("user") != config.RPCUser || h.r.Header.Get("password") != config.RPCPassword {
		msg = "user or password is wrong"
		ok = true
	}
	return
}
func (h *Handler) doHandler() {
	vali, ok := h.Validate()
	if ok {
		h.Err("301", vali)
		return
	}

	body, err := ioutil.ReadAll(h.r.Body)
	if err != nil {

		h.Err("401", "body empty")
		return
	}
	h.SetBody(body)

	if h.r.Header.Get("file") == "upload" {

	} else {

		res, err := Route(h, h.w, h.r)
		if err != nil {
			h.Err(string(res), err.Error())
			return
		}
		h.Out(res)
	}

}
func (h *Handler) SetBody(data []byte) {
	h.body = data
}
func (h *Handler) GetBody() []byte {
	return h.body
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.init(w, r).doHandler()
}
func (h *Handler) RemoteIp() string {
	remoteAddr := h.r.RemoteAddr
	if ip := h.r.Header.Get("XRealIP"); ip != "" {
		remoteAddr = ip
	} else if ip = h.r.Header.Get("XForwardedFor"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

type Bind struct {
	beego.Controller
}

func (i *Bind) Index() {

	if config.RpcServer {
		handler := &Handler{}
		handler.ServeHTTP(i.Ctx.ResponseWriter, i.Ctx.Request)
	}
	return
}
