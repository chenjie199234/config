package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync/atomic"
	"time"

	"github.com/chenjie199234/config/api"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/host"
	"github.com/chenjie199234/Corelib/web"
)

var status int32

//Long Polling
//configservicehost format [http/https]://[username[:password]@]the.host.name[:port]
func NewServiceSdk(selfgroup, selfname, configServiceGroup, configServiceHost string) error {
	if !atomic.CompareAndSwapInt32(&status, 0, 1) {
		return nil
	}
	webclient, e := web.NewWebClient(&web.ClientConfig{
		HeartProbe: time.Second * 5,
	}, selfgroup, selfname, configServiceGroup, api.Name, configServiceHost)
	if e != nil {
		return e
	}
	client := api.NewConfigWebClient(webclient)
	ctx := log.InitTrace(context.Background(), "", selfgroup+"."+selfname, host.Hostip, "watch", "remoteconfig", 0)
	resp, e := client.Watch(ctx, &api.WatchReq{
		Groupname:  selfgroup,
		Appname:    selfname,
		CurVersion: -1,
	}, nil)
	if e != nil {
		return e
	}
	if e = updateApp(common.Str2byte(resp.AppConfig)); e != nil {
		return e
	}
	if e = updateSource(common.Str2byte(resp.SourceConfig)); e != nil {
		return e
	}
	go func() {
		for {
			if e != nil {
				//retry
				time.Sleep(time.Millisecond * 100)
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
			if e = updateApp(common.Str2byte(tmpresp.AppConfig)); e != nil {
				log.Error(nil, "[config.sdk.watch] update appconfig error:", e)
				continue
			}
			if e = updateSource(common.Str2byte(tmpresp.SourceConfig)); e != nil {
				log.Error(nil, "[config.sdk.watch] update sourceconfig error:", e)
				continue
			}
			resp = tmpresp
		}
	}()
	return nil
}
func updateApp(appconfig []byte) error {
	appconfig = bytes.TrimSpace(appconfig)
	if len(appconfig) == 0 {
		appconfig = []byte{'{', '}'}
	}
	if len(appconfig) < 2 || appconfig[0] != '{' || appconfig[len(appconfig)-1] != '}' || !json.Valid(appconfig) {
		return errors.New("[config.sdk.updateApp] data format error")
	}
	appfile, e := os.OpenFile("./AppConfig_tmp.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if e != nil {
		log.Error(nil, "[config.sdk.updateApp] open tmp file error:", e)
		return e
	}
	_, e = appfile.Write(appconfig)
	if e != nil {
		log.Error(nil, "[config.sdk.updateApp] write temp file error:", e)
		return e
	}
	if e = appfile.Sync(); e != nil {
		log.Error(nil, "[config.sdk.updateApp] sync tmp file to disk error:", e)
		return e
	}
	if e = appfile.Close(); e != nil {
		log.Error(nil, "[config.sdk.updateApp] close tmp file error:", e)
		return e
	}
	if e = os.Rename("./AppConfig_tmp.json", "./AppConfig.json"); e != nil {
		log.Error(nil, "[config.sdk.updateApp] rename error:", e)
		return e
	}
	return nil
}
func updateSource(sourceconfig []byte) error {
	sourceconfig = bytes.TrimSpace(sourceconfig)
	if len(sourceconfig) == 0 {
		sourceconfig = []byte{'{', '}'}
	}
	if len(sourceconfig) < 2 || sourceconfig[0] != '{' || sourceconfig[len(sourceconfig)-1] != '}' || !json.Valid(sourceconfig) {
		return errors.New("[config.sdk.updateSource] data format error")
	}
	sourcefile, e := os.OpenFile("./SourceConfig_tmp.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if e != nil {
		log.Error(nil, "[config.sdk.updateSource] open tmp file error:", e)
		return e
	}
	_, e = sourcefile.Write(sourceconfig)
	if e != nil {
		log.Error(nil, "[config.sdk.updateSource] write tmp file error:", e)
		return e
	}
	if e = sourcefile.Sync(); e != nil {
		log.Error(nil, "[config.sdk.updateSource] sync tmp file to disk error:", e)
		return e
	}
	if e = sourcefile.Close(); e != nil {
		log.Error(nil, "[config.sdk.updateSource] close tmp file error:", e)
		return e
	}
	if e = os.Rename("./SourceConfig_tmp.json", "./SourceConfig.json"); e != nil {
		log.Error(nil, "[config.sdk.updateSource] rename error:", e)
		return e
	}
	return nil
}
