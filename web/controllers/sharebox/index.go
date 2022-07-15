package sharebox

import (
	"fmt"
	"io"
	"mmschainnewaccount/config"
	"mmschainnewaccount/sharebox"
	"os"
	"path/filepath"

	"github.com/astaxie/beego"
)

type Index struct {
	beego.Controller
}

func (this *Index) Index() {

	this.TplName = "store/index.tpl"

}

type fileinfoObj struct {
	*sharebox.FileIndex
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

}

func (this *Index) AddFile() {
	out := make(map[string]interface{})

	this.Data["json"] = out
	this.ServeJSON()
	return

}

func (this *Index) GetFile() {
	fn := this.Ctx.Input.Param(":hash")

	fileinfo := sharebox.FindFileindex(fn)
	if fileinfo == nil {

		var err error
		fileinfo, err = sharebox.DownloadFileindexOpt(fn)
		if fileinfo == nil || err != nil {
			fmt.Println("")
			return
		}
	}
	fmt.Println("", fileinfo)

	fileAbsPath := ""

	shareFile := sharebox.FindFile(fileinfo.Hash.B58String())
	if shareFile == nil {

		err := sharebox.DownloadFileOpt(fileinfo)
		if err != nil {
			fmt.Println("", err)
			return
		}
		fmt.Println("")
		fileAbsPath = filepath.Join(config.Store_files, fileinfo.Hash.B58String())
		fileNowPath := filepath.Join(config.Store_files, fileinfo.Name)

		os.Rename(fileAbsPath, fileNowPath)
		fileAbsPath = fileNowPath

	} else {
		fileAbsPath = shareFile.Path
	}

	file, err := os.Open(fileAbsPath)
	if err != nil {

		return
	}

	io.Copy(this.Ctx.ResponseWriter, file)
	file.Close()

	return
}
