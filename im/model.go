package im

import (
	"bytes"
	"fmt"
)

type FriendInfo struct {
	Nickname string `json:"nickname"`
	Hello    string `json:"hello"`
}

func (friendinfo *FriendInfo) Json() []byte {
	res, err := json.Marshal(friendinfo)
	if err != nil {
		fmt.Println(err)
	}
	return res
}
func ParseFriendInfo(bs []byte) *FriendInfo {
	friendinfo := new(FriendInfo)

	decoder := json.NewDecoder(bytes.NewBuffer(bs))
	decoder.UseNumber()
	err := decoder.Decode(friendinfo)
	if err != nil {
		fmt.Println(err)
	}
	return friendinfo
}
