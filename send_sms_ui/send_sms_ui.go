// Send SMS UI
//
// compile for 64bit Linux production using:
// $ env GOOS=linux GOARCH=amd64 go build send_sms_ui.go
//
// Need to add:
// 		config/config.properties file
// 		certs/server.rsa.crt - or whatever is configured in the config file
// 		certs/server.rsa.key - or whatever is configured in the config file
// 		templates/send_sms/all.tmpl - all the template files, if there are more copy them all
//

package main

import (
	"log"
	"os"
	"path/filepath"

	"send_sms_ui/config"
	"send_sms_ui/server"
)

func main() {
	// log.Printf("len(os.Args) = %d\n", len(os.Args))
	// test first argument is used to change certain behaviour e.g. turn off logging to file so can see output in terminal
	test := false
	if len(os.Args) > 1 && os.Args[1] == "test" {
		test = true
	}

	if !test {
		// create log file in same directory as executable
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(filepath.Dir(ex))
		f, err := os.OpenFile(filepath.Dir(ex)+"/send_sms_ui.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	// read config file
	config := config.ReadProperties()

	// run http server
	server.Server(config)
}
