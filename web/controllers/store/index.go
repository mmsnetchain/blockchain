package store

import (
	"fmt"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/libp2parea/virtual_node"
	"io"
	"io/ioutil"
	"mmschainnewaccount/chain_witness_vote/mining/name"
	"mmschainnewaccount/cloud_space"
	"mmschainnewaccount/cloud_space/fs"
	"mmschainnewaccount/config"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/astaxie/beego"
)

type Index struct {
	beego.Controller
}

func (this *Index) Index() {

	this.TplName = "store/index.tpl"

}

type fileinfoObj struct {
	*cloud_space.FileIndex
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

func (this *Index) GetList() {
	_, fileInfos := cloud_space.GetFileindexToSelfAll()

	var result map[string]interface{} = make(map[string]interface{})
	var fileinfoList fileinfoObjList

	for _, fileinfo := range fileInfos {
		var fiobj fileinfoObj
		fiobj.HasCode = fileinfo.Hash.B58String()
		fiobj.FileIndex = fileinfo

		fileinfoList = append(fileinfoList, &fiobj)

	}
	sort.Sort(fileinfoList)
	result["data"] = fileinfoList
	result["status"] = 200
	this.Data["json"] = &result
	this.ServeJSON()
}

func (this *Index) AddFile() {
	fmt.Println("-----")
	out := make(map[string]interface{})

	filePath := this.Ctx.Request.FormValue("file")
	fmt.Println("filepath:", filePath)

	_, filename := filepath.Split(filePath)

	fmt.Println("----- 1111111111")

	fi, err := cloud_space.Diced(filePath)
	fi.Name = filename
	if err != nil {
		fmt.Println("", err)

		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}
	fmt.Println("----- 2222222222222")

	vnodeinfo := nodeStore.NodeSelf.IdInfo.Id

	fi.AddFileOwner(vnodeinfo)

	err = cloud_space.AddFileindexToSelf(fi)
	if err != nil {
		fmt.Println("", err)
		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}
	fmt.Println("----- 333333333333333")

	fmt.Println("----- 4444444444444444")

	out["Code"] = 0
	out["Size"] = fi.Size
	out["HashName"] = fi.Hash.B58String()
	this.Data["json"] = out
	this.ServeJSON()
	return

}

func (this *Index) AddCryptFile() {
	fmt.Println("-----")
	out := make(map[string]interface{})

	nameinfo := name.FindNameToNet(nodeStore.NodeSelf.IdInfo.Id.B58String())

	if nameinfo == nil {
		fmt.Println("")
		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}
	if nameinfo.Deposit < cloud_space.DepositMin {
		fmt.Println("")
		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}

	hs, err := this.GetFiles("files[]")
	if err != nil {
		fmt.Println("", err)
		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}

	f, err := hs[0].Open()
	defer f.Close()
	if err != nil {
		fmt.Println("", err)
		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}
	filename := hs[0].Filename

	newfile, err := os.OpenFile(filepath.Join(config.Store_temp, filename), os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer newfile.Close()

	st, _ := newfile.Stat()
	vnodeinfo := fs.GetNotUseSpace(uint64(st.Size()))
	if vnodeinfo == nil {

		fmt.Println("", err)
		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}
	if err != nil {
		fmt.Println("", err)
		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}

	if _, err := io.Copy(newfile, f); err != nil {
		fmt.Println("", err)

		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}
	data, err := ioutil.ReadFile(filepath.Join(config.Store_temp, filename))
	if err != nil {
		fmt.Println("")
		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}
	dc, err := cloud_space.Encrypt(data)
	if err != nil {
		fmt.Println("")
		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}

	cfilename := filename + cloud_space.Suffix
	cryptfile, err := os.Create(filepath.Join(config.Store_fileinfo_cache, cfilename))
	if err != nil {
		fmt.Println(err)
		return
	}
	cryptfile.Write(dc)
	cryptfile.Close()

	fi, err := cloud_space.Diced(filepath.Join(config.Store_fileinfo_cache, cfilename))
	fi.CryptUser = &nodeStore.NodeSelf.IdInfo.Id
	fi.Name = cfilename
	if err != nil {
		fmt.Println("", err)

		out["Code"] = 1
		out["Message"] = ""
		this.Data["json"] = out
		this.ServeJSON()
		return
	}

	out["Code"] = 0
	out["Size"] = fi.Size
	out["HashName"] = fi.Hash.B58String()
	this.Data["json"] = out
	this.ServeJSON()
	return

}

func (this *Index) GetFile() {
	fn := this.Ctx.Input.Param(":hash")

	fid := virtual_node.AddressFromB58String(fn)

	dw := this.Ctx.Input.Param(":down")
	var isDown bool
	if dw == "down" {
		isDown = true
	} else {
		isDown = false
	}
	haveLocal := true

	fileinfo := cloud_space.FindFileindex(fid)
	if fileinfo == nil {
		haveLocal = false
		var err error
		fileinfo, err = cloud_space.FindFileindexOpt(fn)
		if fileinfo == nil || err != nil {
			fmt.Println("")
			return
		}
	}
	fmt.Printf(" %+v", fileinfo)
	if fileinfo.CryptUser != nil && fileinfo.CryptUser.B58String() != nodeStore.NodeSelf.IdInfo.Id.B58String() {
		fmt.Println("")
		return
	}

	err := cloud_space.DownloadFileOpt(fileinfo, isDown)
	if err != nil {

		return
	}

	vnodes := virtual_node.GetVnodeSelf()
	if len(vnodes) <= 0 {
		return
	}
	if !haveLocal {

		cloud_space.AddFileindexToLocal(fileinfo, vnodes[0].Vid)
	}
	filename := fileinfo.Name

	if fileinfo.CryptUser != nil && fileinfo.CryptUser.B58String() == nodeStore.NodeSelf.IdInfo.Id.B58String() {
		data, err := ioutil.ReadFile(filepath.Join(config.Store_temp, fileinfo.Name))
		if err != nil {
			fmt.Println("")
			return
		}
		dc, err := cloud_space.Decrypt(data)
		if err != nil {
			fmt.Println("")
			return
		}

		filename = strings.TrimRight(filename, cloud_space.Suffix)
		cryptfile, err := os.Create(filepath.Join(config.Store_temp, filename))
		if err != nil {
			fmt.Println(err)
			return
		}
		cryptfile.Write(dc)
		cryptfile.Close()
		os.Remove(filepath.Join(config.Store_temp, fileinfo.Name))
	}

	file, err := os.Open(filepath.Join(config.Store_temp, filename))
	if err != nil {

		return
	}
	io.Copy(this.Ctx.ResponseWriter, file)
	file.Close()

	return

}
