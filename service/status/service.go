package status

import (
	"context"
	"time"

	//"github.com/chenjie199234/config/config"
	//"github.com/chenjie199234/config/ecode"
	"github.com/chenjie199234/config/api"
	statusdao "github.com/chenjie199234/config/dao/status"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/log"
	//"github.com/chenjie199234/Corelib/web"
)

//Service subservice for status business
type Service struct {
	statusDao *statusdao.Dao
}

//Start -
func Start() *Service {
	return &Service{
		//statusDao: statusdao.NewDao(config.GetSql("status_sql"), config.GetRedis("status_redis"), config.GetMongo("status_mongo")),
		statusDao: statusdao.NewDao(nil, nil, nil),
	}
}

func (s *Service) Ping(ctx context.Context, in *api.Pingreq) (*api.Pingresp, error) {
	//if _, ok := ctx.(*crpc.Context); ok {
	//        log.Info("this is a crpc call")
	//}
	//if _, ok := ctx.(*cgrpc.Context); ok {
	//        log.Info("this is a cgrpc call")
	//}
	//if _, ok := ctx.(*web.Context); ok {
	//        log.Info("this is a web call")
	//}
	return &api.Pingresp{ClientTimestamp: in.Timestamp, ServerTimestamp: time.Now().UnixNano()}, nil
}

//Stop -
func (s *Service) Stop() {

}
