package sharebox

import (
	"mmschainnewaccount/rpc/model"
	"mmschainnewaccount/sharebox"
	"net/http"
)

func ShareFolderList(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	rootDir := sharebox.GetShareFolderRootsDetail()
	res, err = model.Tojson(rootDir)
	return
}

func AddLocalShareFoler(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	pathItr, ok := rj.Get("path")
	if !ok {
		res, err = model.Errcode(5002, "path")
		return
	}
	absPath := pathItr.(string)

	err = sharebox.AddLocalShareFolders(absPath)
	if err == nil {
		res, err = model.Tojson("success")
	}
	return
}

func DelLocalShareFoler(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	pathItr, ok := rj.Get("path")
	if !ok {
		res, err = model.Errcode(5002, "path")
		return
	}
	absPath := pathItr.(string)

	err = sharebox.DelLocalShareFolders(absPath)
	if err == nil {
		res, err = model.Tojson("success")
	}
	return
}

func GetRemoteShareFolderList(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	idItr, ok := rj.Get("id")
	if !ok {
		res, err = model.Errcode(5002, "id")
		return
	}
	id := idItr.(string)

	var rootDir *sharebox.DirVO
	rootDir, err = sharebox.GetRemoteShareFolderDetail(id)
	if err != nil {

		res, err = model.Errcode(model.Nomarl, "fail")
		return
	}

	res, err = model.Tojson(rootDir)
	return
}
