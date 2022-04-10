package xcrpc

import (
	"strings"
	"time"

	"github.com/chenjie199234/config/api"
	"github.com/chenjie199234/config/config"
	"github.com/chenjie199234/config/service"

	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/crpc/mids"
	"github.com/chenjie199234/Corelib/log"
	ctime "github.com/chenjie199234/Corelib/util/time"
)

var s *crpc.CrpcServer

//StartCrpcServer -
func StartCrpcServer() {
	c := config.GetCrpcServerConfig()
	crpcc := &crpc.ServerConfig{
		ConnectTimeout: time.Duration(c.ConnectTimeout),
		GlobalTimeout:  time.Duration(c.GlobalTimeout),
		HeartPorbe:     time.Duration(c.HeartProbe),
	}
	var e error
	if s, e = crpc.NewCrpcServer(crpcc, api.Group, api.Name); e != nil {
		log.Error(nil, "[xcrpc] new error:", e)
		return
	}
	UpdateHandlerTimeout(config.AC.HandlerTimeout)

	//this place can register global midwares
	//s.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusCrpcServer(s, service.SvcStatus, mids.AllMids())
	api.RegisterConfigCrpcServer(s, service.SvcConfig, mids.AllMids())
	//example
	//api.RegisterExampleCrpcServer(s, service.SvcExample,mids.AllMids())

	if e = s.StartCrpcServer(":9000"); e != nil && e != crpc.ErrServerClosed {
		log.Error(nil, "[xcrpc] start error:", e)
		return
	}
	log.Info(nil, "[xcrpc] server closed")
}

//UpdateHandlerTimeout -
//first key path,second key method,value timeout duration
func UpdateHandlerTimeout(hts map[string]map[string]ctime.Duration) {
	if s == nil {
		return
	}
	cc := make(map[string]time.Duration)
	for path, methods := range hts {
		for method, timeout := range methods {
			method = strings.ToUpper(method)
			if method == "CRPC" {
				cc[path] = timeout.StdDuration()
			}
		}
	}
	s.UpdateHandlerTimeout(cc)
}

//StopCrpcServer -
func StopCrpcServer() {
	if s != nil {
		s.StopCrpcServer()
	}
}
