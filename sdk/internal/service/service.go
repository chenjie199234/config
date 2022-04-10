package service

import (
	"context"
	"time"

	"github.com/chenjie199234/config/api"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/trace"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/host"
	"github.com/chenjie199234/Corelib/web"
)

type updater func([]byte) error

//serverhost format [http/https]://[username[:password]@]the.host.name[:port]
func NewServiceSdk(selfgroup, selfname, configServiceGroup, serverhost string, updateapp, updatesource updater) error {
	webclient, e := web.NewWebClient(&web.ClientConfig{
		HeartProbe: time.Second * 5,
	}, selfgroup, selfname, configServiceGroup, api.Name, serverhost)
	if e != nil {
		return e
	}
	client := api.NewConfigWebClient(webclient)
	ctx := trace.InitTrace(context.Background(), "", selfgroup+"."+selfname, host.Hostip, "watch", "remoteconfig", 0)
	resp, e := client.Watch(ctx, &api.WatchReq{
		Groupname:  selfgroup,
		Appname:    selfname,
		CurVersion: -1,
	}, nil)
	if e != nil {
		return e
	}
	if e = updateapp(common.Str2byte(resp.AppConfig)); e != nil {
		return e
	}
	if e = updatesource(common.Str2byte(resp.SourceConfig)); e != nil {
		return e
	}
	go func() {
		for {
			if e != nil {
				//retry
				time.Sleep(time.Millisecond * 5)
			}
			tmpresp, e := client.Watch(ctx, &api.WatchReq{
				Groupname:  selfgroup,
				Appname:    selfname,
				CurVersion: resp.Version,
			}, nil)
			if e != nil {
				log.Error(nil, "[config.sdk.watch] error:", e)
				continue
			}
			if e = updateapp(common.Str2byte(tmpresp.AppConfig)); e != nil {
				log.Error(nil, "[config.sdk.watch] update appconfig error:", e)
				continue
			}
			if e = updatesource(common.Str2byte(tmpresp.SourceConfig)); e != nil {
				log.Error(nil, "[config.sdk.watch] update sourceconfig error:", e)
				continue
			}
			resp = tmpresp
		}
	}()
	return nil
}
