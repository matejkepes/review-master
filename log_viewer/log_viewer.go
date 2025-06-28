// Log Viewer - For viewing log files in a web browser
//
// compile for 64bit Linux production using:
// $ env GOOS=linux GOENV=amd64 go build log_viewer.go
//
// Need to add:
// 		config/config.properties file
// 		certs/server.rsa.crt - or whatever is configured in the config file
// 		certs/server.rsa.key - or whatever is configured in the config file
//

package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"

	"log_viewer/config"
	"log_viewer/server"
)

func main() {
	// log.Printf("len(os.Args) = %d\n", len(os.Args))
	// test first argument is used to change certain behaviour e.g. turn off logging to file so can see output in terminal
	test := false
	if len(os.Args) > 1 && os.Args[1] == "test" {
		test = true
	}

	logFilename := ""
	if !test {
		// create log file in same directory as executable
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(filepath.Dir(ex))
		logFilename = filepath.Dir(ex) + "/log_viewer.log"
		// f, err := os.OpenFile(filepath.Dir(ex)+"/log_viewer.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		f, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		// log.SetOutput(f)
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		})
	}

	// read config file
	config := config.ReadProperties()

	// run http server
	server.Server(config)
}
