package cmd

import "log"

// Application struct is a wrap of application, which controls everything and have top management role, everything is united arround this struct.
type Application struct {
	Debug    *bool
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
