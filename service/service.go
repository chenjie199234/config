package service

import (
	"github.com/chenjie199234/config/dao"
	"github.com/chenjie199234/config/service/config"
	"github.com/chenjie199234/config/service/status"
)

//SvcStatus one specify sub service
var SvcStatus *status.Service

//SvcConfig one specify sub service
var SvcConfig *config.Service

//StartService start the whole service
func StartService() error {
	if e := dao.NewApi(); e != nil {
		return e
	}
	//start sub service
	SvcStatus = status.Start()
	SvcConfig = config.Start()
	return nil
}

//StopService stop the whole service
func StopService() {
	//stop sub service
	SvcStatus.Stop()
	SvcConfig.Stop()
}
