package sdk

import (
	"encoding/json"
	"errors"
	"os"
	"sync/atomic"

	"github.com/chenjie199234/config/sdk/internal/direct"
	"github.com/chenjie199234/config/sdk/internal/service"

	"github.com/chenjie199234/Corelib/log"
)

var status int32

//url format [mongodb/mongodb+srv]://[username:password@]host1,...,hostN[/dbname][?param1=value1&...&paramN=valueN]
func NewDirectSdk(selfgroup, selfname string, url string) error {
	if atomic.CompareAndSwapInt32(&status, 0, 1) {
		return direct.NewDirectSdk(selfgroup, selfname, url, updateApp, updateSource)
	}
	return nil
}

//configservicehost format [http/https]://[username[:password]@]the.host.name[:port]
func NewServiceSdk(selfgroup, selfname, configservicegroup, configservicehost string) error {
	if atomic.CompareAndSwapInt32(&status, 0, 2) {
		return service.NewServiceSdk(selfgroup, selfname, configservicegroup, configservicehost, updateApp, updateSource)
	}
	return nil
}
func updateApp(appconfig []byte) error {
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
