package mining

import (
	"errors"
)

var (
	ERROR_repeat_import_block                = errors.New("Repeat import block")
	ERROR_fork_import_block                  = errors.New("The front block cannot be found, and the new block is discontinuous")
	ERROR_import_block_height_not_continuity = errors.New("Import block height discontinuity")
)
