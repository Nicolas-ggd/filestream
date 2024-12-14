package app

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	APP_PORT = os.Getenv("APP_PORT")
)

func App() {
	// parse addr flag to dynamicaly change address value
	addr := flag.String("addr", APP_PORT, "HTTP network address")

	// parse debug mode to enable debug mode in application
	debug := flag.Bool("debug", false, "Enable debug mode")
	// define logs
	infoLog := log.New(os.Stderr, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Llongfile)

	flag.Parse()

	if *debug {
		infoLog.Println("Enabled debug mode...")
	}

	_ = &Application{
		Debug:    debug,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	srv := http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatalln(err)
	}
}
