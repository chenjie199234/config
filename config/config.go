package config

import (
	"os"
	"strconv"

	"github.com/chenjie199234/config/api"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/trace"
	"github.com/chenjie199234/config/sdk"
)

//EnvConfig can't hot update,all these data is from system env setting
//nil field means that system env not exist
type EnvConfig struct {
	ConfigType *int
	RunEnv     *string
	DeployEnv  *string
}

//EC -
var EC *EnvConfig

//notice is a sync function
//don't write block logic inside it
func Init(notice func(c *AppConfig)) {
	initenv()
	initremote()
	initsource()
	initapp(notice)
}

//Close -
func Close() {
	log.Close()
	trace.Close()
}

func initenv() {
	EC = &EnvConfig{}
	if str, ok := os.LookupEnv("CONFIG_TYPE"); ok && str != "<CONFIG_TYPE>" && str != "" {
		configtype, e := strconv.Atoi(str)
		if e != nil || (configtype != 0 && configtype != 1) {
			log.Error(nil, "[config.initenv] env CONFIG_TYPE must be number in [0,1]")
			Close()
			os.Exit(1)
		}
		EC.ConfigType = &configtype
	} else {
		log.Warning(nil, "[config.initenv] missing env CONFIG_TYPE")
	}
	if str, ok := os.LookupEnv("RUN_ENV"); ok && str != "<RUN_ENV>" && str != "" {
		EC.RunEnv = &str
	} else {
		log.Warning(nil, "[config.initenv] missing env RUN_ENV")
	}
	if str, ok := os.LookupEnv("DEPLOY_ENV"); ok && str != "<DEPLOY_ENV>" && str != "" {
		EC.DeployEnv = &str
	} else {
		log.Warning(nil, "[config.initenv] missing env DEPLOY_ENV")
	}
}

func initremote() {
	if EC.ConfigType == nil || *EC.ConfigType == 0 {
		return
	}
	if *EC.ConfigType == 1 {
		var mongourl string
		if str, ok := os.LookupEnv("REMOTE_CONFIG_MONGO_URL"); ok && str != "<REMOTE_CONFIG_MONGO_URL>" && str != "" {
			mongourl = str
		} else {
			panic("[config.initremote] missing env REMOTE_CONFIG_MONGO_URL")
		}
		if e := sdk.NewDirectSdk(api.Group, api.Name, mongourl); e != nil {
			log.Error(nil, "[config.initremote] new sdk error:", e)
			Close()
			os.Exit(1)
		}
	}
}
