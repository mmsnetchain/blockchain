package rpc

import (
	"github.com/prestonTao/libp2parea/virtual_node"
	"mmschainnewaccount/cloud_space/fs"
	"mmschainnewaccount/rpc/model"
	"net/http"
)

func SetSpaceSize(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	nItr, ok := rj.Get("n")
	if !ok {
		res, err = model.Errcode(5002, "n")
		return
	}
	n := uint64(nItr.(float64))
	virtual_node.SetupVnodeNumber(n)
	res, err = model.Tojson("success")
	return
}

func GetSpaceSize(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	total := fs.GetSpaceSize()
	useSize := fs.GetUseSpaceSize()
	result := make(map[string]uint64, 0)
	result["TotalSize"] = total
	result["UseSize"] = useSize

	res, err = model.Tojson(result)
	return
}

func AddSpaceSize(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	nItr, ok := rj.Get("n")
	if !ok {
		res, err = model.Errcode(model.NoField, "n")
		return
	}
	n := uint64(nItr.(float64))

	absPath := ""

	absPathItr, ok := rj.Get("absPath")
	if ok {
		absPath = absPathItr.(string)
	}

	fs.AddSpace(absPath, n)

	res, err = model.Tojson("success")
	return
}

func DelSpaceSize(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	nItr, ok := rj.Get("n")
	if !ok {
		res, err = model.Errcode(model.NoField, "n")
		return
	}
	n := uint64(nItr.(float64))

	fs.DelSpace(n)

	res, err = model.Tojson("success")
	return
}
