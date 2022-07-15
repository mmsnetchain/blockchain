package im

import (
	"bytes"
	"fmt"
)

type FileInfo struct {
	Name  string
	Size  int64
	Path  string
	Index int64
	Data  []byte
}

func (fi *FileInfo) Json() []byte {
	d, err := json.Marshal(fi)
	if err != nil {
		fmt.Println(err)
	}
	return d
}
func ParseFileInfo(d []byte) *FileInfo {
	fi := FileInfo{}

	decoder := json.NewDecoder(bytes.NewBuffer(d))
	decoder.UseNumber()
	err := decoder.Decode(&fi)
	if err != nil {
		fmt.Println(err)
	}
	return &fi
}
