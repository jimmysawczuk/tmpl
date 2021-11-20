package tmplfunc

import "os"

type Depender interface {
	Depend(string) error
}

type Filesystem interface {
	In() *os.File
	Out() *os.File
	BaseDir() string
}

type FilesystemDepender interface {
	Depender
	Filesystem
}
