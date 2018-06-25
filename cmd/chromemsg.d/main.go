package main

import (
	"fmt"
	"os"
	"time"

	// internal - core
	chromemsg "github.com/sniperkit/snk.golang.chrome-extension/pkg"
)

func main() {
	fmt.Println("Running chrome message cli command...")

	chromeMsgService, err := chromemsg.NewWithConfig(
		&chromemsg.Config{
			PortType: chromemsg.Bufio,
			Debug:    true,
		},
	)
	if err != nil {
		panic(err)
	}

	chromeMsgHello := map[string]interface{}{
		"pkg":  "chromemsg",
		"time": time.Now().Format("2015-06-16-0431 UTC"),
	}

	// chromeMsgService.Read()
	if err := chromeMsgService.Write(chromeMsgHello); err != nil {
		panic(err)
	}

	os.Exit(1)
}
