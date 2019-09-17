package common

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func SignalWaitQuit(functor *Functor) {
	var stopLock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		//阻塞程序运行，直到收到终止的信号
		s := <-signalChan
		stopLock.Lock()
		stop = true
		stopLock.Unlock()
		LogInfo("catch signal: %v", s)
		functor.Call(s)
		stopChan <- struct{}{}
		os.Exit(0)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
}
