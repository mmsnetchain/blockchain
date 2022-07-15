package db

func SaveHistory(id []byte, bs *[]byte) error {
	return LevelTempDB.Save(id, bs)

}
