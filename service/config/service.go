package config

import (
	"context"
	"encoding/json"

	"github.com/chenjie199234/config/api"
	"github.com/chenjie199234/config/config"
	configdao "github.com/chenjie199234/config/dao/config"
	"github.com/chenjie199234/config/ecode"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/common"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
)

//Service subservice for config business
type Service struct {
	configDao *configdao.Dao
}

//Start -
func Start() *Service {
	return &Service{
		configDao: configdao.NewDao(config.GetSql("config_sql"), config.GetRedis("config_redis"), config.GetMongo("config_mongo")),
	}
}

//get all groups
func (s *Service) Groups(ctx context.Context, req *api.GroupsReq) (*api.GroupsResp, error) {
	groups, e := s.configDao.MongoGetAllGroups(ctx, req.SearchFilter)
	if e != nil {
		log.Error(ctx, "[Groups]", e)
		return nil, ecode.ErrSystem
	}
	return &api.GroupsResp{Groups: groups}, nil
}

//get all apps
func (s *Service) Apps(ctx context.Context, req *api.AppsReq) (*api.AppsResp, error) {
	apps, e := s.configDao.MongoGetAllApps(ctx, req.Groupname, req.SearchFilter)
	if e != nil {
		log.Error(ctx, "[Apps] group:", req.Groupname, "error:", e)
		return nil, ecode.ErrSystem
	}
	return &api.AppsResp{Apps: apps}, nil
}

//get one specific app's config
func (s *Service) Get(ctx context.Context, req *api.GetReq) (*api.GetResp, error) {
	summary, config, e := s.configDao.MongoGetCOnfig(ctx, req.Groupname, req.Appname, req.Index)
	if e != nil {
		log.Error(ctx, "[Get] group:", req.Groupname, "app:", req.Appname, "error:", e)
		return nil, ecode.ErrSystem
	}
	if summary == nil {
		return &api.GetResp{}, nil
	}
	if config == nil {
		if req.Index == 0 {
			log.Error(ctx, "[Get] group:", req.Groupname, "app:", req.Appname, "error: has summary but missing current config")
			return nil, ecode.ErrSystem
		}
		log.Error(ctx, "[Get] group:", req.Groupname, "app:", req.Appname, "index:", req.Index, "error: doesn't exist")
		return nil, ecode.ErrReq
	}
	return &api.GetResp{
		CurIndex:     summary.CurIndex,
		MaxIndex:     summary.MaxIndex,
		CurVersion:   summary.CurVersion,
		ThisIndex:    config.Index,
		AppConfig:    config.AppConfig,
		SourceConfig: config.SourceConfig,
	}, nil
}

//set one specific app's config
func (s *Service) Set(ctx context.Context, req *api.SetReq) (*api.SetResp, error) {
	if req.AppConfig == "" {
		req.AppConfig = "{}"
	} else if len(req.AppConfig) < 2 || req.AppConfig[0] != '{' || req.AppConfig[len(req.AppConfig)-1] != '}' || !json.Valid(common.Str2byte(req.AppConfig)) {
		log.Error(ctx, "[Set] group:", req.Groupname, "app:", req.Appname, "AppConfig data not in json object format")
		return nil, ecode.ErrReq
	}
	if req.SourceConfig == "" {
		req.SourceConfig = "{}"
	} else if len(req.SourceConfig) < 2 || req.SourceConfig[0] != '{' || req.SourceConfig[len(req.SourceConfig)-1] != '}' || !json.Valid(common.Str2byte(req.SourceConfig)) {
		log.Error(ctx, "[Set] group:", req.Groupname, "app:", req.Appname, "SourceConfig data not in json object format")
		return nil, ecode.ErrReq
	}
	if e := s.configDao.MongoSetConfig(ctx, req.Groupname, req.Appname, req.AppConfig, req.SourceConfig); e != nil {
		log.Error(ctx, "[Set] group:", req.Groupname, "app:", req.Appname, "error:", e)
		return nil, ecode.ErrSystem
	}
	return &api.SetResp{}, nil
}

//rollback one specific app's config
func (s *Service) Rollback(ctx context.Context, req *api.RollbackReq) (*api.RollbackResp, error) {
	status, e := s.configDao.MongoRollbackConfig(ctx, req.Groupname, req.Appname, req.Index)
	if e != nil {
		log.Error(ctx, "[Rollback] group:", req.Groupname, "app:", req.Appname, "error:e", e)
		return nil, ecode.ErrSystem
	}
	if !status {
		log.Error(ctx, "[Rollback] group:", req.Groupname, "app:", req.Appname, "index:", req.Index, "doesn't exist")
		return nil, ecode.ErrReq
	}
	return &api.RollbackResp{}, nil
}

//Stop -
func (s *Service) Stop() {

}
