package config

import (
	"path/filepath"
)

const (
	SQLITE3DB_name               = "sqlite3.db"
	Table_name_friend            = "friend"
	Table_name_shareFolderRemote = "share_folder_remote"
	SQL_SHOW                     = false
)

var (
	SQLITE3DB_path = filepath.Join(Path_configDir, SQLITE3DB_name)
)
