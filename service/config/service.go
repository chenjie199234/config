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
	config  *model.Current
	notices map[chan *struct{}]*struct{}
}

//Start -
func Start() *Service {
	s := &Service{
		configDao:  configdao.NewDao(config.GetSql("config_sql"), config.GetRedis("config_redis"), config.GetMongo("config_mongo")),
		noticepool: &sync.Pool{},
		apps:       make(map[string]map[string]*app),
		status:     true,
	}
	if e := s.configDao.MongoWatchConfig(s.refresh, s.update, s.delgroup, s.delapp, s.delconfig); e != nil {
		panic("[Config.Start] watch error: " + e.Error())
	}
	return s
}

//first key groupname,second key appname,value curconfig
func (s *Service) refresh(curs map[string]map[string]*model.Current) {
	s.Lock()
	defer s.Unlock()
	//delete not exist
	for gname, g := range s.apps {
		curg, ok := curs[gname]
		if !ok {
			log.Info(nil, "[refresh.delgroup] group:", gname)
			for _, a := range g {
				if len(a.notices) == 0 {
					//if there are no watchers,clean right now
					delete(g, a.config.AppName)
				} else {
					//if there are watchers,clean will happened when watcher return
					a.config.CurVersion = 0
					a.config.AppConfig = "{}"
					a.config.SourceConfig = "{}"
					for notice := range a.notices {
						notice <- nil
					}
				}
			}
			if len(g) == 0 {
				//if there are no watchers,clean right now
				delete(s.apps, gname)
			}
			continue
		}
		for _, a := range g {
			if _, ok := curg[a.config.AppName]; !ok {
				log.Info(nil, "[refresh.delapp] group:", a.config.GroupName, "app:", a.config.AppName)
				if len(a.notices) == 0 {
					//if there are no watchers,clean right now
					delete(g, a.config.AppName)
				} else {
					//if there are watchers,clean will happened when watcher return
					a.config.CurVersion = 0
					a.config.AppConfig = "{}"
					a.config.SourceConfig = "{}"
					for notice := range a.notices {
						notice <- nil
					}
				}
			}
		}
		if len(g) == 0 {
			//if there are no watchers,clean right now
			delete(s.apps, gname)
		}
	}
	//add new or refresh exist
	for gname, curg := range curs {
		g, ok := s.apps[gname]
		if !ok {
			g = make(map[string]*app)
			s.apps[gname] = g
		}
		for aname, cura := range curg {
			log.Info(nil, "[refresh.update] group:", gname, "app:", aname, "version:", cura.CurVersion, "AppConfig:", cura.AppConfig, "SourceConfig:", cura.SourceConfig)
			a, ok := g[aname]
			if !ok {
				//this is a new
				if cura.CurVersion == 0 {
					//curversion 0,this is same as doesn't exist this app
					continue
				}
				a = &app{
					config:  cura,
					notices: make(map[chan *struct{}]*struct{}),
				}
				g[aname] = a
				continue
			}
			//already exist
			if cura.CurVersion == 0 && len(a.notices) == 0 {
				//curversion 0,this is same as doesn't exist this app
				//if there are no watchers,clean right now
				delete(g, aname)
				continue
			}
			a.config = cura
			for notice := range a.notices {
				notice <- nil
			}
		}
		if len(g) == 0 {
			delete(s.apps, gname)
		}
	}
}
func (s *Service) update(cur *model.Current) {
	log.Info(nil, "[update] group:", cur.GroupName, "app:", cur.AppName, "version:", cur.CurVersion, "AppConfig:", cur.AppConfig, "SourceConfig:", cur.SourceConfig)
	s.Lock()
	defer s.Unlock()
	g, ok := s.apps[cur.GroupName]
	if !ok {
		g = make(map[string]*app)
		s.apps[cur.GroupName] = g
	}
	a, ok := g[cur.AppName]
	if !ok {
		//this is a new
		if cur.CurVersion == 0 {
			//curversion 0,this is same as doesn't exist this app
			if len(g) == 0 {
				delete(s.apps, cur.GroupName)
			}
			return
		}
		a = &app{
			config:  cur,
			notices: make(map[chan *struct{}]*struct{}),
		}
		g[cur.AppName] = a
		return
	}
	//already exist
	if cur.CurVersion == 0 && len(a.notices) == 0 {
		//curversion 0,this is same as doesn't exist this app
		//if there are no watchers,clean right now
		delete(g, a.config.AppName)
		if len(g) == 0 {
			delete(s.apps, cur.GroupName)
		}
		return
	}
	a.config = cur
	for notice := range a.notices {
		notice <- nil
	}
}
func (s *Service) delgroup(groupname string) {
	s.Lock()
	defer s.Unlock()
	g, ok := s.apps[groupname]
	if !ok {
		return
	}
	log.Info(nil, "[delgroup] group:", groupname)
	for _, a := range g {
		if len(a.notices) == 0 {
			//if there are no watchers,clean right now
			delete(g, a.config.AppName)
		} else {
			//if there are watchers,clean will happened when watcher return
			a.config.CurVersion = 0
			a.config.AppConfig = "{}"
			a.config.SourceConfig = "{}"
			for notice := range a.notices {
				notice <- nil
			}
		}
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
	log.Info(nil, "[delapp] group:", groupname, "app:", appname)
	if len(a.notices) == 0 {
		//if there are no watchers,clean right now
		delete(g, a.config.AppName)
		if len(g) == 0 {
			delete(s.apps, groupname)
		}
	} else {
		//if there are watchers,clean will happened when watcher return
		a.config.CurVersion = 0
		a.config.AppConfig = "{}"
		a.config.SourceConfig = "{}"
		for notice := range a.notices {
			notice <- nil
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
	if a.config.SummaryID != summaryid {
		return
	}
	//delete the summary,this is same as delete the app
	log.Info(nil, "[delconfig] group:", groupname, "app:", appname)
	if len(a.notices) == 0 {
		//if there are no watchers,clean right now
		delete(g, a.config.AppName)
		if len(g) == 0 {
			delete(s.apps, groupname)
		}
	} else {
		//if there are watchers,clean will happened when watcher return
		a.config.CurVersion = 0
		a.config.AppConfig = "{}"
		a.config.SourceConfig = "{}"
		for notice := range a.notices {
			notice <- nil
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
	//when MongoGetConfig's param index is 0
	//summary and config must be both nil or both not nil
	if summary == nil {
		return &api.GetResp{}, nil
	}
	if config == nil {
		//this can only happened when the index is not 0
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
	if int32(a.config.CurVersion) != req.CurVersion {
		resp := &api.WatchResp{
			Version:      int32(a.config.CurVersion),
			AppConfig:    a.config.AppConfig,
			SourceConfig: a.config.SourceConfig,
		}
		s.Unlock()
		return resp, nil
	}
	ch := s.getnotice()
	a.notices[ch] = nil
	s.Unlock()
	select {
	case <-ctx.Done():
		s.Lock()
		defer s.Unlock()
		delete(a.notices, ch)
		s.putnotice(ch)
		if len(a.notices) == 0 && a.config.CurVersion == 0 {
			delete(g, a.config.AppName)
		}
		if len(g) == 0 {
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
	delete(a.notices, ch)
	s.putnotice(ch)
	if len(a.notices) == 0 && a.config.CurVersion == 0 {
		delete(g, a.config.AppName)
	}
	if len(g) == 0 {
		delete(s.apps, a.config.GroupName)
	}
	return &api.WatchResp{
		Version:      int32(a.config.CurVersion),
		AppConfig:    a.config.AppConfig,
		SourceConfig: a.config.SourceConfig,
	}, nil
}

//Stop -
func (s *Service) Stop() {
	s.Lock()
	defer s.Unlock()
	s.status = false
	for _, g := range s.apps {
		for _, a := range g {
			for n := range a.notices {
				close(n)
			}
		}
	}
}
