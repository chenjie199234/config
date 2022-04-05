package config

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/chenjie199234/config/api"
	"github.com/chenjie199234/config/config"
	configdao "github.com/chenjie199234/config/dao/config"
	"github.com/chenjie199234/config/ecode"

	cerror "github.com/chenjie199234/Corelib/error"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/common"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
)

//Service subservice for config business
type Service struct {
	configDao  *configdao.Dao
	noticepool *sync.Pool
	sync.Mutex
	groups map[string]*group
}
type group struct {
	sync.Mutex
	apps map[string]*app
}
type app struct {
	sync.Mutex
	config  *configdata
	notices map[chan *struct{}]*struct{}
}
type configdata struct {
	Version      int32
	AppConfig    string
	SoucreConfig string
}

//Start -
func Start() *Service {
	s := &Service{
		configDao:  configdao.NewDao(config.GetSql("config_sql"), config.GetRedis("config_redis"), config.GetMongo("config_mongo")),
		noticepool: &sync.Pool{},
	}
	// if e := s.update(); e != nil {
	// panic("")
	// }
	return s
}

//
// func (s *Service) update() error {
// }

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

func (s *Service) getnotice() chan *struct{} {
	ch, ok := s.noticepool.Get().(chan *struct{})
	if !ok {
		return make(chan *struct{}, 1)
	}
	return ch
}
func (s *Service) putnotice(ch chan *struct{}) {
	s.noticepool.Put(ch)
}

//watch on specific app's config
func (s *Service) Watch(ctx context.Context, req *api.WatchReq) (*api.WatchResp, error) {
	s.Lock()
	g, ok := s.groups[req.Groupname]
	if !ok {
		g = &group{}
		s.groups[req.Groupname] = g
	}
	g.Lock()
	s.Unlock()
	a, ok := g.apps[req.Appname]
	if !ok {
		a = &app{
			config: &configdata{
				Version:      0,
				AppConfig:    "{}",
				SoucreConfig: "{}",
			},
			notices: make(map[chan *struct{}]*struct{}),
		}
		g.apps[req.Appname] = a
	}
	a.Lock()
	g.Unlock()
	if req.CurVersion < 0 || a.config.Version > req.CurVersion {
		resp := &api.WatchResp{
			Version:      a.config.Version,
			AppConfig:    a.config.AppConfig,
			SourceConfig: a.config.SoucreConfig,
		}
		a.Unlock()
		return resp, nil
	} else if a.config.Version < req.CurVersion {
		curversion := a.config.Version
		a.Unlock()
		log.Error(ctx, "[Watch] client version:", req.CurVersion, "big then current version:", curversion)
		return nil, ecode.ErrReq
	} else {
		ch := s.getnotice()
		a.notices[ch] = nil
		a.Unlock()
		if _, ok := <-ch; ok {
			s.putnotice(ch)
		} else {
			return nil, cerror.ErrClosing
		}
	}
	a.Lock()
	resp := &api.WatchResp{
		Version:      a.config.Version,
		AppConfig:    a.config.AppConfig,
		SourceConfig: a.config.SoucreConfig,
	}
	a.Unlock()
	return resp, nil
}

//Stop -
func (s *Service) Stop() {
	s.Lock()
	defer s.Unlock()
	for _, g := range s.groups {
		g.Lock()
		for _, a := range g.apps {
			a.Lock()
			for n := range a.notices {
				close(n)
			}
			a.Unlock()
		}
		g.Unlock()
	}
}
