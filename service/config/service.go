package config

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/chenjie199234/config/api"
	"github.com/chenjie199234/config/config"
	configdao "github.com/chenjie199234/config/dao/config"
	"github.com/chenjie199234/config/ecode"
	"github.com/chenjie199234/config/model"

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
	apps   map[string]map[string]*app //first key groupname,second key appname,value appinfo
	status bool
}
type app struct {
	sync.Mutex
	config  *model.Current
	notices map[chan *struct{}]*struct{}
}

//Start -
func Start() *Service {
	s := &Service{
		configDao:  configdao.NewDao(config.GetSql("config_sql"), config.GetRedis("config_redis"), config.GetMongo("config_mongo")),
		noticepool: &sync.Pool{},
		status:     true,
	}
	if e := s.configDao.MongoWatchConfig(s.refresh, s.update, s.delgroup, s.delapp, s.delconfig); e != nil {
		panic("[Config.Start] watch error: " + e.Error())
	}
	return s
}
func (s *Service) refresh(curs []*model.Current) {

}
func (s *Service) update(cur *model.Current) {

}
func (s *Service) delgroup(groupname string) {
	s.Lock()
	defer s.Unlock()
	g, ok := s.apps[groupname]
	if !ok {
		return
	}
	for _, a := range g {
		a.Lock()
		a.config.CurVersion = 0
		a.config.AppConfig = "{}"
		a.config.SourceConfig = "{}"
		for notice := range a.notices {
			notice <- nil
		}
		if len(a.notices) == 0 {
			//if there are no watchers,clean right now
			delete(g, a.config.AppName)
		}
		a.Unlock()
	}
	if len(g) == 0 {
		//if there are no watchers,clean right now
		delete(s.apps, groupname)
	}
}
func (s *Service) delapp(groupname, appname string) {
	s.Lock()
	defer s.Unlock()
	g, ok := s.apps[groupname]
	if !ok {
		return
	}
	a, ok := g[appname]
	if !ok {
		return
	}
	a.Lock()
	defer a.Unlock()
	a.config.CurVersion = 0
	a.config.AppConfig = "{}"
	a.config.SourceConfig = "{}"
	for notice := range a.notices {
		notice <- nil
	}
	if len(a.notices) == 0 {
		//if there are no watchers,clean right now
		delete(g, a.config.AppName)
		if len(g) == 0 {
			delete(s.apps, groupname)
		}
	}
}
func (s *Service) delconfig(groupname, appname, summaryid string) {
	s.Lock()
	defer s.Unlock()
	g, ok := s.apps[groupname]
	if !ok {
		return
	}
	a, ok := g[appname]
	if !ok {
		return
	}
	a.Lock()
	defer a.Unlock()
	if a.config.SummaryID != summaryid {
		return
	}
	//delete the summary,this is same as delete the app
	a.config.CurVersion = 0
	a.config.AppConfig = "{}"
	a.config.SourceConfig = "{}"
	for notice := range a.notices {
		notice <- nil
	}
	if len(a.notices) == 0 {
		//if there are no watchers,clean right now
		delete(g, a.config.AppName)
		if len(g) == 0 {
			delete(s.apps, groupname)
		}
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
	summary, config, e := s.configDao.MongoGetConfig(ctx, req.Groupname, req.Appname, req.Index)
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
	if !s.status {
		s.Unlock()
		return nil, cerror.ErrClosing
	}
	g, ok := s.apps[req.Groupname]
	if !ok {
		g = make(map[string]*app)
		s.apps[req.Groupname] = g
	}
	a, ok := g[req.Appname]
	if !ok {
		a = &app{
			config: &model.Current{
				CurVersion:   0,
				GroupName:    req.Groupname,
				AppName:      req.Appname,
				AppConfig:    "{}",
				SourceConfig: "{}",
			},
			notices: make(map[chan *struct{}]*struct{}),
		}
		g[req.Appname] = a
	}
	a.Lock()
	s.Unlock()
	if int32(a.config.CurVersion) != req.CurVersion {
		resp := &api.WatchResp{
			Version:      int32(a.config.CurVersion),
			AppConfig:    a.config.AppConfig,
			SourceConfig: a.config.SourceConfig,
		}
		a.Unlock()
		return resp, nil
	}
	ch := s.getnotice()
	a.notices[ch] = nil
	a.Unlock()
	select {
	case <-ctx.Done():
		s.Lock()
		defer s.Unlock()
		a.Lock()
		defer a.Unlock()
		delete(a.notices, ch)
		s.putnotice(ch)
		if len(a.notices) == 0 && a.config.CurVersion == 0 {
			delete(s.apps[a.config.GroupName], a.config.AppName)
		}
		if len(s.apps[a.config.GroupName]) == 0 {
			delete(s.apps, a.config.GroupName)
		}
		return nil, ctx.Err()
	case _, ok := <-ch:
		if !ok {
			return nil, cerror.ErrClosing
		}
	}
	s.Lock()
	defer s.Unlock()
	a.Lock()
	defer a.Unlock()
	delete(a.notices, ch)
	s.putnotice(ch)
	resp := &api.WatchResp{
		Version:      int32(a.config.CurVersion),
		AppConfig:    a.config.AppConfig,
		SourceConfig: a.config.SourceConfig,
	}
	if len(a.notices) == 0 && a.config.CurVersion == 0 {
		delete(s.apps[a.config.GroupName], a.config.AppName)
	}
	if len(s.apps[a.config.GroupName]) == 0 {
		delete(s.apps, a.config.GroupName)
	}
	return resp, nil
}

//Stop -
func (s *Service) Stop() {
	s.Lock()
	defer s.Unlock()
	s.status = false
	for _, g := range s.apps {
		for _, a := range g {
			a.Lock()
			for n := range a.notices {
				close(n)
			}
			a.Unlock()
		}
	}
}
