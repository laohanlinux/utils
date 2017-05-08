// +build linux

package platform

import (
	"os"
	"os/signal"

	log "github.com/laohanlinux/utils/gokitlog"
)

type SigFn func()

var sigFns []SigFn

func RegisterSignal(sig ...os.Signal) {
	signalChan := make(chan os.Signal)
	go func() {
		for {
			log.Info("receive the term signal", <-signalChan)
			for _, fn := range sigFns {
				fn()
			}
		}
	}()
	signal.Notify(signalChan, sig...)
}

func RegisterSigFn(fns ...SigFn) {
	sigFns = fns
}
