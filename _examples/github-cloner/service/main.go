package main

import (
	"fmt"
	"os"

	// external
	"github.com/k0kubun/pp"
)

func main() {
	fmt.Println("Running chrome message cli command...")

	// default logger

	fmt.Println("logger backend=", defaultLoggerConfig.Backend.String())
	log = newLoggerWithConfig(nil)

	cmsg, err := initChromeMsg(nil)
	fatalOnError(err)

	pp.Println("cmsg=", cmsg)

	// default service
	svc := newService(nil)

	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	} else {
		fmt.Println("you need to provide at least one command...")
		os.Exit(1)
	}

	status, err := svc.Manage(command)
	fatalOnError(err)
	fmt.Println(status)

}

func fatalOnError(err error) {
	if err != nil {
		fmt.Println("error: ", err.Error())
		os.Exit(1)
	}
}
