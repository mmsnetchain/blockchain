package im

import (
	"mmschainnewaccount/sqlite3_db"
)

func AddFriend_opt() {}

func GetContactsList() []sqlite3_db.Friends {
	fs, err := new(sqlite3_db.Friends).Getall()
	if err != nil {
		return nil
	}
	return fs
}

func DelContacts(id string) {
	new(sqlite3_db.Friends).Del(id)
	new(sqlite3_db.MsgLog).RemoveAllForFriend(id)
}
