package main

import (
	"time"

	// internal - core
	chromemsg "github.com/sniperkit/snk.golang.chromemsg/pkg"
)

// defaultChromeMsgConfig sets the default chrome native messenger configuration struct
var defaultChromeMsgConfig = &chromemsg.Config{
	PortType: chromemsg.Bufio,
	Debug:    true,
}

// initChromeMsg instanciates a new chrome native messenger (read/write)...
func initChromeMsg(conf *chromemsg.Config) (*chromemsg.Messenger, error) {
	if conf == nil {
		conf = defaultChromeMsgConfig
	}

	chromeMsgService, err := chromemsg.NewWithConfig(conf)
	if err != nil {
		panic(err)
	}

	chromeMsgHello := map[string]interface{}{
		"pkg":  "chromemsg",
		"time": time.Now().Format("2015-06-16-0431 UTC"),
	}

	// chromeMsgService.Read()
	if err := chromeMsgService.Write(chromeMsgHello); err != nil {
		return nil, err
	}
	return chromeMsgService, nil
}
