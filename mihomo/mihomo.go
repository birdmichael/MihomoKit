package mihomo

import (
	"fmt"

	"github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/log"
)

type RealTimeLogger interface {
	Log(level string, payload string)
}

var logger RealTimeLogger

func Setup(homeDir string, config string) {
	if logger == nil {
		fmt.Println("No logger set. Mihomo will not log.")
	}
	go SubscribeLogger()
	constant.SetHomeDir(homeDir)
	constant.SetHomeDir(config)
}

func SetupLogger(l RealTimeLogger) {
	logger = l
}

func SubscribeLogger() {
	sub := log.Subscribe()
	defer log.UnSubscribe(sub)

	for ev := range sub {
		if logger != nil {
			log := ev
			logger.Log(log.Type(), log.Payload)
		}
	}
}