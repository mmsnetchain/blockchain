package im

import (
	"github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/config"
	"path/filepath"
)

func RegisterIM() {
	message_center.Register_p2pHE(config.MSGID_im_file, FileMsg)
	message_center.Register_p2pHE(config.MSGID_im_file_recv, FileMsg_recv)
	message_center.Register_p2pHE(config.MSGID_im_property, PropertyMsg)
	message_center.Register_p2pHE(config.MSGID_im_property_recv, PropertyMsg_recv)
	message_center.Register_p2pHE(config.MSGID_im_addr, BaseCoinAddrMsg)
	message_center.Register_p2pHE(config.MSGID_im_addr_recv, BaseCoinAddrMsg_recv)
	message_center.Register_p2pHE(config.MSGID_im_pay, PayMsg)
	message_center.Register_p2pHE(config.MSGID_im_pay_recv, PayMsg_recv)

	utils.CheckCreateDir(filepath.Join(imfilepath))
}
