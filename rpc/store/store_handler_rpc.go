package store

import (
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/libp2parea/virtual_node"
	"mmschainnewaccount/rpc/model"

	"mmschainnewaccount/sharebox"
	ystore "mmschainnewaccount/store"

	"fmt"
	"mmschainnewaccount/store/fs"
	"net/http"
	"sort"
)

func UploadFile() {
	fmt.Println("123456789")
}

type fileinfoObj struct {
	*ystore.FileIndex
	HasCode string
}
type fileinfoObjList []*fileinfoObj

func (fil fileinfoObjList) Len() int {
	return len(fil)
}
func (fil fileinfoObjList) Less(i, j int) bool {
	return fil[i].FileIndex.Time > fil[j].FileIndex.Time
}
func (fil fileinfoObjList) Swap(i, j int) {
	fil[i], fil[j] = fil[j], fil[i]
}

func GetFileList(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	_, fileInfos := ystore.GetFileindexToSelfAll()

	var fileinfoList fileinfoObjList

	for _, fileinfo := range fileInfos {
		var fiobj fileinfoObj
		fiobj.HasCode = fileinfo.Hash.B58String()
		fiobj.FileIndex = fileinfo

		fileinfoList = append(fileinfoList, &fiobj)

	}
	sort.Sort(fileinfoList)
	res, err = model.Tojson(fileinfoList)
	return
}

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
		return
	}

	res, err = model.Tojson(rootDir)
	return
}

func SearchFileInfo(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	hash, ok := rj.Get("hash")
	if !ok {
		res, err = model.Errcode(5002, "hash")
		return
	}
	fn := hash.(string)

	fid := virtual_node.AddressFromB58String(fn)

	fileindex := ystore.FindFileindex(fid)
	if fileindex == nil {
		fileindex, err = ystore.FindFileindexOpt(fn)
		if fileindex == nil || err != nil {
			fmt.Println("")
			res, err = model.Errcode(5003, "")
			return
		}
	}
	if fileindex.CryptUser != nil && fileindex.CryptUser.B58String() != nodeStore.NodeSelf.IdInfo.Id.B58String() {
		fmt.Println("")
		res, err = model.Errcode(5003, "")
		return
	}
	res, err = model.Tojson(fileindex)
	return
}

func DelFileInfo(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	hash, ok := rj.Get("hash")
	if !ok {
		res, err = model.Errcode(5002, "hash")
		return
	}
	fn := hash.(string)
	fid := virtual_node.AddressFromB58String(fn)

	err = ystore.DelFileInfoFromSelf(fid)
	if err != nil {
		fmt.Println("", err.Error())
		res, err = model.Errcode(model.Nomarl, err.Error())
	} else {
		res, err = model.Tojson("")
	}
	return
}

func AddFileInfo(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	hash, ok := rj.Get("hash")
	if !ok {
		res, err = model.Errcode(5002, "hash")
		return
	}
	fn := hash.(string)
	fid := virtual_node.AddressFromB58String(fn)

	finfo := ystore.FindFileindex(fid)
	if finfo == nil {
		finfo, err = ystore.FindFileindexOpt(fn)
		if finfo == nil || err != nil {
			fmt.Println("")
			res, err = model.Errcode(5002, "file not fund")
			return
		}
	}
	fmt.Printf(":%+v", finfo)
	if finfo.CryptUser != nil && finfo.CryptUser.B58String() != nodeStore.NodeSelf.IdInfo.Id.B58String() {
		fmt.Println("")
		res, err = model.Errcode(5002, "")
		return
	}

	vnodeinfo := fs.GetNotUseSpace(finfo.Size)

	finfo.AddFileOwner(*vnodeinfo)

	err = ystore.AddFileindexToSelf(finfo, vnodeinfo.Vid)
	res, err = model.Tojson("")
	return
}

func DownloadProcOne(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	hash, ok := rj.Get("hash")
	if !ok {
		res, err = model.Errcode(5002, "hash")
		return
	}
	fn := hash.(string)

	vid := virtual_node.AddressFromB58String(fn)

	finfo := ystore.FindFileindex(vid)
	dp := ystore.DownloadProgressOne(finfo)
	res, err = model.Tojson(dp)
	return
}

func DownloadProc(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	dp := ystore.DownloadProgress()
	res, err = model.Tojson(dp)
	return
}

func DownloadComplete(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	dp := ystore.DownLoadComplete()
	res, err = model.Tojson(dp)
	return
}

func DownLoadStop(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	hash, ok := rj.Get("hash")
	if !ok {
		res, err = model.Errcode(5002, "hash")
		return
	}
	fn := hash.(string)
	ystore.DownLoadStop(fn)
	res, err = model.Tojson("")
	return
}

func DownLoadDel(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	hash, ok := rj.Get("hash")
	if !ok {
		res, err = model.Errcode(5002, "hash")
		return
	}
	fn := hash.(string)
	ystore.DwonLoadDel(fn)
	res, err = model.Tojson("")
	return
}

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

type VnodeinfoVO struct {
	Nid   string `json:"nid"`
	Index uint64 `json:"index"`
	Vid   string `json:"vid"`
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
