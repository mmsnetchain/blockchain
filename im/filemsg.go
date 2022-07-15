package im

import (
	"bytes"
	"fmt"
	"github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/libp2parea/message_center/flood"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"io/ioutil"
	"mmschainnewaccount/config"
	"mmschainnewaccount/sqlite3_db"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	Count  int64 = 20 * 1024 * 1024 * 1024
	Lenth  int64 = 200 * 1024
	ErrNum int   = 5
	Second int64 = 10
)

type Msg struct {
	From     nodeStore.AddressNet
	To       nodeStore.AddressNet
	Nowid    int64
	Toid     int64
	Text     []byte
	FilePath string
	File     *FileInfo
	Class    int
	Speed    map[string]int64
}

func (msg *Msg) SetSpeed(stime int64, size int) error {

	if _, ok := msg.Speed["time"]; !ok {
		msg.Speed["time"] = stime
		msg.Speed["size"] = int64(size)
	}
	if time.Now().Unix()-msg.Speed["time"] > Second {
		msg.Speed["time"] = stime
		msg.Speed["size"] = 0
	} else {
		msg.Speed["size"] += int64(size)
	}
	return nil
}

func (msg *Msg) GetSpeed() int64 {
	t := time.Now().Unix() - msg.Speed["time"]
	if t < 1 {
		t = 1
	}
	return msg.Speed["size"] / t * 100
}
func (msg *Msg) Json() []byte {
	d, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	return d
}

func ParseMsg(d []byte) (*Msg, error) {
	msg := &Msg{}

	decoder := json.NewDecoder(bytes.NewBuffer(d))
	decoder.UseNumber()
	err := decoder.Decode(msg)
	if err != nil {
		fmt.Println(err)
	}
	return msg, err
}

func (msg *Msg) ReadFile() []byte {
	stat, err := os.Stat(msg.FilePath)
	if err != nil {
		fmt.Println(err)
	}
	d, err := ioutil.ReadFile(msg.FilePath)
	if err != nil {
		fmt.Println(err)
	}
	fi := FileInfo{Name: stat.Name(), Size: stat.Size(), Data: d}
	msg.File = &fi
	return msg.Json()
}

func (msg *Msg) ReadFileSlice(id int64) (fd []byte, fileinfo FileInfo, ok bool, errs error) {

	msginfo, err := new(sqlite3_db.MsgLog).FindById(id)
	if err != nil {
		fmt.Println(err)
		errs = err
		return
	}

	if msginfo.Index >= msginfo.Size && msginfo.Size > 0 {
		return
	}

	msg.To = nodeStore.AddressFromB58String(msginfo.Recipient)
	msg.Text = []byte(msginfo.Content)
	msg.FilePath = msginfo.Path
	msg.Class = msginfo.Class
	msg.Nowid = id
	msg.Toid = msginfo.Toid
	path := msginfo.Path
	index := msginfo.Index
	fi, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		errs = err
		return
	}
	stat, err := fi.Stat()
	if err != nil {
		fmt.Println(err)
		errs = err
		return
	}
	size := stat.Size()
	start := index
	length := Lenth

	if start < size && size-index < Lenth {
		length = size - index
	}
	buf := make([]byte, length)
	_, err = fi.ReadAt(buf, start)
	if err != nil {
		fmt.Println(err)
		errs = err
		return
	}

	nextstart := start + Lenth
	if nextstart >= size {
		ok = true
	}
	fmt.Println("...", size, nextstart)
	finfo := FileInfo{Name: stat.Name(), Size: size, Index: nextstart, Data: buf}

	msg.File = &finfo
	fd = msg.Json()
	fileinfo = finfo
	return
}

func (msg *Msg) SendFile(id int64) (bl bool) {

	var errnum int
	for i := int64(0); ; i++ {
	BEGIN:
		bs, fi, okf, err := msg.ReadFileSlice(id)
		if bs == nil {
			break
		}
		if err != nil {

			errnum++
			if errnum <= ErrNum {
				fmt.Println("resend slice...")
				goto BEGIN
			}
			break
		}
		message, ok, _ := message_center.SendP2pMsgHE(config.MSGID_im_file, &msg.To, &bs)
		if ok {

			rbs := flood.WaitRequest(config.CLASS_im_file_msg, utils.Bytes2string(message.Body.Hash), 0)
			if rbs != nil {
				toid := utils.BytesToInt64(*rbs)

				rate := int64(float64(fi.Index) / float64(fi.Size) * float64(100))
				msg.SetSpeed(time.Now().Unix(), len(bs))
				speed := msg.GetSpeed()

				new(sqlite3_db.MsgLog).UpRate(id, toid, fi.Index, rate, speed, fi.Size)
				if okf {
					bl = true
					break
				}
			} else {

				fmt.Println("fail")
				errnum++
				if errnum <= ErrNum {

					fmt.Println("resend...")
					goto BEGIN
				}
				bl = false
				break
			}
		} else {
			bl = false
			break
		}
	}
	return
}
