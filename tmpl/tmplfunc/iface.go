package tmplfunc

import "os"

type Refer interface {
	Ref(string) error
}

type Filesystem interface {
	In() *os.File
	Out() *os.File
	BaseDir() string
}

type FilesystemRefer interface {
	Refer
	Filesystem
}

type Moder interface {
	IsProduction() bool
}
