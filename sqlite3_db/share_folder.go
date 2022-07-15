package sqlite3_db

const (
	split = "*|*"
)

type ShareFolder struct {
	Path string `xorm:"varchar(25) pk notnull unique 'path'"`
}

func (this *ShareFolder) Add(path string) error {

	this.Path = path
	_, err := engineDB.Insert(this)
	return err
}

func (this *ShareFolder) Del(path string) {

	this.Path = path
	engineDB.Where("path=?", path).Unscoped().Delete(this)
}

func (this *ShareFolder) GetAll() ([]ShareFolder, error) {
	sf := make([]ShareFolder, 0)
	err := engineDB.Find(&sf)
	return sf, err
}
