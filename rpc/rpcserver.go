package rpc

import (
	"fmt"
	"mmschainnewaccount/rpc/model"
	"net/http"
)

func parseJson(jsonb []byte) (*model.RpcJson, error) {
	var rpcjson model.RpcJson
	err := json.Unmarshal(jsonb, &rpcjson)

	return &rpcjson, err
}
func Route(rh model.RpcHandler, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	data := rh.GetBody()
	rj, err := parseJson(data)
	if err != nil {

	}
	hd, ok := rpcHandler[rj.Method]
	if ok {
		res, err = hd(rj, w, r)
	} else {
		res, err = model.Errcode(model.NoMethod, rj.Method)
	}
	return
}

func UploadFile(rh model.RpcHandler) (res []byte, err error) {
	fmt.Println("")

	return nil, nil
}
