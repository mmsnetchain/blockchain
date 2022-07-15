package store

import (
	"mmschainnewaccount/rpc/model"
	"mmschainnewaccount/store"
	"net/http"
)

func AddFolder(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	name, ok := rj.Get("name")
	if !ok {
		res, err = model.Errcode(5002, "name")
		return
	}
	fname := name.(string)
	var pid uint64
	parentid, ok := rj.Get("parentid")
	if ok {
		pid = uint64(parentid.(float64))
	}
	err = store.AddFolder(pid, fname)
	if err == nil {
		res, err = model.Tojson("success")
	}
	return
}

func DelFolder(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	idstr, ok := rj.Get("id")
	if !ok {
		res, err = model.Errcode(5002, "id")
		return
	}
	id := idstr.(float64)
	err = store.DelFolder(uint64(id))
	if err == nil {
		res, err = model.Tojson("success")
	}
	return
}

func UpFolder(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	idstr, ok := rj.Get("id")
	if !ok {
		res, err = model.Errcode(5002, "id")
		return
	}
	id := uint64(idstr.(float64))
	namestr, ok := rj.Get("name")
	if !ok {
		res, err = model.Errcode(5002, "name")
		return
	}
	name := namestr.(string)
	var pid uint64
	parentid, ok := rj.Get("parentid")
	if ok {
		pid = uint64(parentid.(float64))
	}
	err = store.UpFolder(id, pid, name)
	if err != nil {
		res, err = model.Tojson(err.Error())
		return
	}
	res, err = model.Tojson("success")
	return
}

func ListFolder(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	pidstr, ok := rj.Get("parentid")
	if !ok {
		res, err = model.Errcode(5002, "parentid")
		return
	}
	pid := uint64(pidstr.(float64))
	list := store.ListFolder(pid)
	res, err = model.Tojson(list)
	return
}

func Moveto(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	hashstr, ok := rj.Get("hash")
	if !ok {
		res, err = model.Errcode(5002, "hash")
		return
	}
	hash := hashstr.(string)
	pidstr, ok := rj.Get("pid")
	if !ok {
		res, err = model.Errcode(5002, "parentid")
		return
	}
	pid := uint64(pidstr.(float64))
	b := store.Moveto(hash, pid)
	res, err = model.Tojson(b)
	return
}
