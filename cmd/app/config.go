package app

import (
	"github.com/Nicolas-ggd/filestream"
	"log"
)

// Application struct is a wrap of application, which controls everything and have top management role, everything is united around this struct.
type Application struct {
	Debug        *bool
	ErrorLog     *log.Logger
	InfoLog      *log.Logger
	RFileRequest *fstream.RFileRequest
	File         *fstream.File
}
