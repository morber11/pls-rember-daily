package filemode

import "os"

const (
	OwnerReadWriteExecute os.FileMode = 0o700
	OwnerReadWrite        os.FileMode = 0o600
)
