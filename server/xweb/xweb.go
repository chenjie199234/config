package xweb

import (
	"net/http"
	"strings"
	"time"

	"github.com/chenjie199234/config/api"
	"github.com/chenjie199234/config/config"
	"github.com/chenjie199234/config/service"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/web/mids"
)

var s *web.WebServer

//StartWebServer -
func StartWebServer() {
	c := config.GetWebServerConfig()
	webc := &web.ServerConfig{
		ConnectTimeout:     time.Duration(c.ConnectTimeout),
		GlobalTimeout:      time.Duration(c.GlobalTimeout),
		IdleTimeout:        time.Duration(c.IdleTimeout),
		HeartProbe:         time.Duration(c.HeartProbe),
		StaticFileRootPath: c.StaticFilePath,
		MaxHeader:          1024,
	}
	if c.Cors != nil {
		webc.Cors = &web.CorsConfig{
			AllowedOrigin:    c.Cors.CorsOrigin,
			AllowedHeader:    c.Cors.CorsHeader,
			ExposeHeader:     c.Cors.CorsExpose,
			AllowCredentials: true,
			MaxAge:           24 * time.Hour,
		}
	}
	var e error
	if s, e = web.NewWebServer(webc, api.Group, api.Name); e != nil {
		log.Error(nil, "[xweb] new error:", e)
		return
	}
	UpdateHandlerTimeout(config.AC)

	//this place can register global midwares
	//s.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusWebServer(s, service.SvcStatus, mids.AllMids())
	api.RegisterConfigWebServer(s, service.SvcConfig, mids.AllMids())
	//example
	//api.RegisterExampleWebServer(s, service.SvcExample, mids.AllMids())

	if e = s.StartWebServer(":8000"); e != nil && e != web.ErrServerClosed {
		log.Error(nil, "[xweb] start error:", e)
		return
	}
	log.Info(nil, "[xweb] server closed")
}

//UpdateHandlerTimeout -
func UpdateHandlerTimeout(c *config.AppConfig) {
	if s != nil {
		cc := make(map[string]map[string]time.Duration)
		for path, methods := range c.HandlerTimeout {
			for method, timeout := range methods {
				if timeout == 0 {
					continue
				}
				method = strings.ToUpper(method)
				if method != http.MethodGet && method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch && method != http.MethodDelete {
					continue
				}
				if _, ok := cc[method]; !ok {
					cc[method] = make(map[string]time.Duration)
				}
				cc[method][path] = timeout.StdDuration()
			}
		}
		s.UpdateHandlerTimeout(cc)
	}
}

//StopWebServer -
func StopWebServer() {
	if s != nil {
		s.StopWebServer()
	}
}