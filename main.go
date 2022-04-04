package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/chenjie199234/config/config"
	"github.com/chenjie199234/config/server/xcrpc"
	"github.com/chenjie199234/config/server/xgrpc"
	"github.com/chenjie199234/config/server/xweb"
	"github.com/chenjie199234/config/service"

	"github.com/chenjie199234/Corelib/log"
	publicmids "github.com/chenjie199234/Corelib/mids"
	_ "github.com/chenjie199234/Corelib/monitor"
)

func main() {
	config.Init(func(ac *config.AppConfig) {
		//this is a notice callback every time appconfig changes
		//this function works in sync mode
		//don't write block logic inside this
		log.Info(nil, "[main] new app config:", ac)
		xcrpc.UpdateHandlerTimeout(ac)
		xgrpc.UpdateHandlerTimeout(ac)
		xweb.UpdateHandlerTimeout(ac)
		publicmids.UpdateRateConfig(ac.HandlerRate)
		publicmids.UpdateAccessKeyConfig(ac.AccessKeys)
	})
	defer config.Close()
	//start the whole business service
	if e := service.StartService(); e != nil {
		log.Error(nil, e)
		return
	}
	//start low level net service
	ch := make(chan os.Signal, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		xcrpc.StartCrpcServer()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xweb.StartWebServer()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xgrpc.StartCGrpcServer()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	//stop the whole business service
	service.StopService()
	//stop low level net service
	wg.Add(1)
	go func() {
		xcrpc.StopCrpcServer()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xweb.StopWebServer()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xgrpc.StopCGrpcServer()
		wg.Done()
	}()
	wg.Wait()
}
