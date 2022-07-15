package im

import (
	"fmt"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/libp2parea/message_center/flood"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/config"
	"mmschainnewaccount/sqlite3_db"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func FileMsg(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	var id int64

	content := *message.Body.Content
	m, err := ParseMsg(content)
	if err == nil {

		sendId := message.Head.Sender.B58String()
		m.File.Path = filepath.Join(imfilepath, m.File.Name)
		num := 0
	Rename:

		if ok, _ := utils.PathExists(m.File.Path); ok {
			num++
			filenamebase := filepath.Base(m.File.Name)
			fileext := filepath.Ext(m.File.Name)
			filename := strings.TrimSuffix(filenamebase, fileext)
			newname := filename + "_" + strconv.Itoa(num) + fileext
			m.File.Path = filepath.Join(imfilepath, newname)
			if ok1, _ := utils.PathExists(m.File.Path); ok1 {
				goto Rename
			}
			m.File.Name = newname
		}

		tmpPath := filepath.Join(imfilepath, m.File.Name+"_"+sendId+"_tmp")
		fi, err := os.OpenFile(tmpPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		start := m.File.Index - Lenth
		if start >= m.File.Size {
			start = m.File.Size
		}
		fi.Seek(start, 0)
		fi.Write(m.File.Data)
		defer fi.Close()
		fmt.Println(sendId, start)
		ml := sqlite3_db.MsgLog{}
		if m.Toid == 0 {
			id, err = ml.Add(sendId, sqlite3_db.Self, string(m.Text), m.File.Path, m.Class)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			id = m.Toid
		}
		rate := int64(float64(m.File.Index) / float64(m.File.Size) * float64(100))
		m.SetSpeed(time.Now().Unix(), len(content))
		speed := m.GetSpeed()
		err = ml.UpRate(id, m.Nowid, m.File.Index, rate, speed, m.File.Size)
		if err != nil {
			fmt.Println("update transimission rate fail", err)
		}

		if rate >= 100 {

			fi.Close()
			os.Rename(tmpPath, m.File.Path)

			ml.IsSuccessful(id)
			now := time.Now()
			msgVO := core.MessageVO{
				DBID:     id,
				Id:       sendId,
				Index:    now.Unix(),
				Time:     utils.FormatTimeToSecond(now),
				Content:  string(m.Text),
				Path:     m.File.Path,
				FileName: m.File.Name,
				Size:     m.File.Size,
				Cate:     m.Class,
			}
			msgVO.DBID = id
			select {
			case core.MsgChannl <- &msgVO:
			default:
			}
		}
	}

	bs := utils.Int64ToBytes(id)
	message_center.SendP2pReplyMsgHE(message, config.MSGID_im_file_recv, &bs)
}

func FileMsg_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	flood.ResponseWait(config.CLASS_im_file_msg, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}

func PropertyMsg(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	p := new(sqlite3_db.Property).Get(nodeStore.NodeSelf.IdInfo.Id.B58String())

	bs := p.Json()
	message_center.SendP2pReplyMsgHE(message, config.MSGID_im_property_recv, &bs)
}

func PropertyMsg_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	flood.ResponseWait(config.CLASS_im_property_msg, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}

func BaseCoinAddrMsg(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	addr := []byte(keystore.GetCoinbase().Addr)
	message_center.SendP2pReplyMsgHE(message, config.MSGID_im_addr_recv, &addr)
}

func BaseCoinAddrMsg_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	flood.ResponseWait(config.CLASS_im_addr_msg, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}

func PayMsg(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	content := string(*message.Body.Content)

	sendId := message.Head.Sender.B58String()
	ml := sqlite3_db.MsgLog{}
	id, err := ml.Add(sendId, sqlite3_db.Self, content, "", core.MsgPayId)
	if err == nil {
		ml.IsSuccessful(id)
	}

	now := time.Now()
	msgVO := core.MessageVO{
		DBID:    id,
		Id:      sendId,
		Index:   now.Unix(),
		Time:    utils.FormatTimeToSecond(now),
		Content: content,
		Cate:    core.MsgPayId,
	}
	msgVO.DBID = id
	select {
	case core.MsgChannl <- &msgVO:
	default:
	}
	bs := []byte("ok")
	message_center.SendP2pReplyMsgHE(message, config.MSGID_im_pay_recv, &bs)
}

func PayMsg_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	flood.ResponseWait(config.CLASS_im_pay_msg, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}
