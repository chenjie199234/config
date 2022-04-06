package xgrpc

import (
	"strings"
	"time"

	"github.com/chenjie199234/config/api"
	"github.com/chenjie199234/config/config"
	"github.com/chenjie199234/config/service"

	"github.com/chenjie199234/Corelib/cgrpc"
	"github.com/chenjie199234/Corelib/cgrpc/mids"
	"github.com/chenjie199234/Corelib/log"
	ctime "github.com/chenjie199234/Corelib/util/time"
)

var s *cgrpc.CGrpcServer

//StartCGrpcServer -
func StartCGrpcServer() {
	c := config.GetCGrpcServerConfig()
	cgrpcc := &cgrpc.ServerConfig{
		ConnectTimeout: time.Duration(c.ConnectTimeout),
		GlobalTimeout:  time.Duration(c.GlobalTimeout),
		HeartPorbe:     time.Duration(c.HeartProbe),
	}
	var e error
	if s, e = cgrpc.NewCGrpcServer(cgrpcc, api.Group, api.Name); e != nil {
		log.Error(nil, "[xgrpc] new error:", e)
		return
	}
	UpdateHandlerTimeout(config.AC.HandlerTimeout)

	//this place can register global midwares
	//s.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusCGrpcServer(s, service.SvcStatus, mids.AllMids())
	api.RegisterConfigCGrpcServer(s, service.SvcConfig, mids.AllMids())
	//example
	//api.RegisterExampleCGrpcServer(s, service.SvcExample, mids.AllMids())

	if e = s.StartCGrpcServer(":10000"); e != nil && e != cgrpc.ErrServerClosed {
		log.Error(nil, "[xgrpc] start error:", e)
		return
	}
	log.Info(nil, "[xgrpc] server closed")
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
			if timeout == 0 {
				continue
			}
			method = strings.ToUpper(method)
			if method == "GRPC" {
				cc[path] = timeout.StdDuration()
			}
		}
	}
	s.UpdateHandlerTimeout(cc)
}

//StopCGrpcServer -
func StopCGrpcServer() {
	if s != nil {
		s.StopCGrpcServer()
	}
}
