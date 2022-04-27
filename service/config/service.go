package config

import (
	"context"
	"encoding/json"
	"strings"
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
	summary *model.Summary
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
	if e := s.configDao.MongoWatchConfig(s.refresh, s.update, s.delgroup, s.delapp, s.delconfig, model.Decrypt); e != nil {
		panic("[Config.Start] watch error: " + e.Error())
	}
	return s
}

//first key groupname,second key appname,value curconfig
func (s *Service) refresh(curs map[string]map[string]*model.Summary) {
	s.Lock()
	defer s.Unlock()
	//delete not exist
	for gname, g := range s.apps {
		curg, ok := curs[gname]
		if !ok {
			log.Debug(nil, "[refresh.delgroup] group:", gname)
			for aname, a := range g {
				if len(a.notices) == 0 {
					//if there are no watchers,clean right now
					delete(g, aname)
				} else {
					//if there are watchers,clean will happened when watcher return
					a.summary.Cipher = ""
					a.summary.CurVersion = 0
					a.summary.CurAppConfig = "{}"
					a.summary.CurSourceConfig = "{}"
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
		for aname, a := range g {
			if _, ok := curg[aname]; !ok {
				log.Debug(nil, "[refresh.delapp] group:", gname, "app:", aname)
				if len(a.notices) == 0 {
					//if there are no watchers,clean right now
					delete(g, aname)
				} else {
					//if there are watchers,clean will happened when watcher return
					a.summary.Cipher = ""
					a.summary.CurVersion = 0
					a.summary.CurAppConfig = "{}"
					a.summary.CurSourceConfig = "{}"
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
		g, gok := s.apps[gname]
		if !gok {
			g = make(map[string]*app)
		}
		for aname, cura := range curg {
			log.Debug(nil, "[refresh.update] group:", gname, "app:", aname, "version:", cura.CurVersion, "AppConfig:", cura.CurAppConfig, "SourceConfig:", cura.CurSourceConfig)
			a, ok := g[aname]
			if !ok {
				//this is a new
				if cura.CurVersion == 0 {
					//this is same as not exist
					continue
				}
				g[aname] = &app{
					summary: cura,
					notices: make(map[chan *struct{}]*struct{}),
				}
				continue
			}
			//already exist
			if cura.CurVersion == 0 && len(a.notices) == 0 {
				//this is same as not exist and there are no watchers,clean right now
				delete(g, aname)
				continue
			}
			a.summary = cura
			for notice := range a.notices {
				notice <- nil
			}
		}
		if !gok && len(g) > 0 {
			s.apps[gname] = g
		} else if gok && len(g) == 0 {
			delete(s.apps, gname)
		}
	}
}
func (s *Service) update(gname, aname string, cur *model.Summary) {
	log.Debug(nil, "[update] group:", gname, "app:", aname, "version:", cur.CurVersion, "AppConfig:", cur.CurAppConfig, "SourceConfig:", cur.CurSourceConfig)
	s.Lock()
	defer s.Unlock()
	g, gok := s.apps[gname]
	if !gok {
		g = make(map[string]*app)
	}
	defer func() {
		if !gok && len(g) > 0 {
			s.apps[gname] = g
		} else if gok && len(g) == 0 {
			delete(s.apps, gname)
		}
	}()
	a, ok := g[aname]
	if !ok {
		//this is a new
		if cur.CurVersion == 0 {
			//this is same as not exist
			return
		}
		g[aname] = &app{
			summary: cur,
			notices: make(map[chan *struct{}]*struct{}),
		}
		return
	}
	//already exist
	if cur.CurVersion == 0 && len(a.notices) == 0 {
		//this is same as not exist and there are no watchers,clean right now
		delete(g, aname)
		return
	}
	a.summary = cur
	for notice := range a.notices {
		notice <- nil
	}
}
func (s *Service) delgroup(groupname string) {
	log.Debug(nil, "[delgroup] group:", groupname)
	s.Lock()
	defer s.Unlock()
	g, ok := s.apps[groupname]
	if !ok {
		return
	}
	for aname, a := range g {
		if len(a.notices) == 0 {
			//if there are no watchers,clean right now
			delete(g, aname)
		} else {
			//if there are watchers,clean will happened when watcher return
			a.summary.Cipher = ""
			a.summary.CurVersion = 0
			a.summary.CurAppConfig = "{}"
			a.summary.CurSourceConfig = "{}"
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
	log.Debug(nil, "[delapp] group:", groupname, "app:", appname)
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
	if len(a.notices) == 0 {
		//if there are no watchers,clean right now
		delete(g, appname)
		if len(g) == 0 {
			delete(s.apps, groupname)
		}
	} else {
		//if there are watchers,clean will happened when watcher return
		a.summary.Cipher = ""
		a.summary.CurVersion = 0
		a.summary.CurAppConfig = "{}"
		a.summary.CurSourceConfig = "{}"
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
	if a.summary.ID.Hex() != summaryid {
		log.Debug(nil, "[delconfig] group:", groupname, "app:", appname, "config")
		return
	}
	//delete the summary,this is same as delete the app
	log.Debug(nil, "[delconfig] group:", groupname, "app:", appname, "summary")
	if len(a.notices) == 0 {
		//if there are no watchers,clean right now
		delete(g, appname)
		if len(g) == 0 {
			delete(s.apps, groupname)
		}
	} else {
		//if there are watchers,clean will happened when watcher return
		a.summary.Cipher = ""
		a.summary.CurVersion = 0
		a.summary.CurAppConfig = "{}"
		a.summary.CurSourceConfig = "{}"
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

//create one specific app
func (s *Service) Create(ctx context.Context, req *api.CreateReq) (*api.CreateResp, error) {
	if req.Cipher != "" && len(req.Cipher) != 32 {
		log.Error(ctx, "[Create] group:", req.Groupname, "app:", req.Appname, "error:", ecode.ErrCipherLength)
		return nil, ecode.ErrCipherLength
	}
	if e := s.configDao.MongoCreate(ctx, req.Groupname, req.Appname, req.Cipher, model.Encrypt); e != nil {
		log.Error(ctx, "[Create] group:", req.Groupname, "app:", req.Appname, "error:", e)
		if e != ecode.ErrAppAlreadyExist {
			e = ecode.ErrSystem
		}
		return nil, e
	}
	log.Info(ctx, "[Create] group:", req.Groupname, "app:", req.Appname, "success")
	return &api.CreateResp{}, nil
}

//update one specific app's cipher
func (s *Service) Updatecipher(ctx context.Context, req *api.UpdatecipherReq) (*api.UpdatecipherResp, error) {
	if req.New != "" && len(req.New) != 32 {
		log.Error(ctx, "[Updatechiper] group:", req.Groupname, "app:", req.Appname, "error:", ecode.ErrCipherLength)
		return nil, ecode.ErrCipherLength
	}
	if req.Old == req.New {
		return &api.UpdatecipherResp{}, nil
	}
	if e := s.configDao.MongoUpdateCipher(ctx, req.Groupname, req.Appname, req.Old, req.New, model.Decrypt, model.Encrypt); e != nil {
		log.Error(ctx, "[Updatechiper] group:", req.Groupname, "app:", req.Appname, "error:", e)
		if e != ecode.ErrAppNotExist && e != ecode.ErrWrongCipher {
			e = ecode.ErrSystem
		}
		return nil, e
	}
	log.Info(ctx, "[Updatecipher] group:", req.GetGroupname, "app:", req.Appname, "success")
	return &api.UpdatecipherResp{}, nil
}

//get one specific app's config
func (s *Service) Get(ctx context.Context, req *api.GetReq) (*api.GetResp, error) {
	summary, config, e := s.configDao.MongoGetConfig(ctx, req.Groupname, req.Appname, req.Index, model.Decrypt)
	if e != nil {
		log.Error(ctx, "[Get] group:", req.Groupname, "app:", req.Appname, "error:", e)
		return nil, ecode.ErrSystem
	}
	if summary == nil {
		return nil, ecode.ErrAppNotExist
	}
	if config == nil {
		//this can only happened when the index is not 0
		log.Error(ctx, "[Get] group:", req.Groupname, "app:", req.Appname, "index:", req.Index, "error: doesn't exist")
		return nil, ecode.ErrIndexNotExist
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
	req.AppConfig = strings.TrimSpace(req.AppConfig)
	if req.AppConfig == "" {
		req.AppConfig = "{}"
	} else if len(req.AppConfig) < 2 || req.AppConfig[0] != '{' || req.AppConfig[len(req.AppConfig)-1] != '}' || !json.Valid(common.Str2byte(req.AppConfig)) {
		return nil, ecode.ErrConfigFormat
	}
	req.SourceConfig = strings.TrimSpace(req.SourceConfig)
	if req.SourceConfig == "" {
		req.SourceConfig = "{}"
	} else if len(req.SourceConfig) < 2 || req.SourceConfig[0] != '{' || req.SourceConfig[len(req.SourceConfig)-1] != '}' || !json.Valid(common.Str2byte(req.SourceConfig)) {
		return nil, ecode.ErrConfigFormat
	}
	index, e := s.configDao.MongoSetConfig(ctx, req.Groupname, req.Appname, req.AppConfig, req.SourceConfig, model.Encrypt)
	if e != nil {
		log.Error(ctx, "[Set] group:", req.Groupname, "app:", req.Appname, "error:", e)
		if e != ecode.ErrAppNotExist {
			e = ecode.ErrSystem
		}
		return nil, e
	}
	log.Info(ctx, "[Set] group:", req.Groupname, "app:", req.Appname, "index:", index, "success")
	return &api.SetResp{}, nil
}

//rollback one specific app's config
func (s *Service) Rollback(ctx context.Context, req *api.RollbackReq) (*api.RollbackResp, error) {
	if e := s.configDao.MongoRollbackConfig(ctx, req.Groupname, req.Appname, req.Index); e != nil {
		log.Error(ctx, "[Rollback] group:", req.Groupname, "app:", req.Appname, "error:e", e)
		if e != ecode.ErrAppNotExist && e != ecode.ErrIndexNotExist {
			e = ecode.ErrSystem
		}
		return nil, e
	}
	log.Info(ctx, "[Rollback] group:", req.Groupname, "app:", req.Appname, "index:", req.Index, "success")
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
			summary: &model.Summary{
				Cipher:          "",
				CurVersion:      0,
				CurAppConfig:    "{}",
				CurSourceConfig: "{}",
			},
			notices: make(map[chan *struct{}]*struct{}),
		}
		g[req.Appname] = a
	}
	if int32(a.summary.CurVersion) != req.CurVersion {
		resp := &api.WatchResp{
			Version:      int32(a.summary.CurVersion),
			AppConfig:    a.summary.CurAppConfig,
			SourceConfig: a.summary.CurSourceConfig,
		}
		s.Unlock()
		return resp, nil
	}
	for {
		ch := s.getnotice()
		a.notices[ch] = nil
		s.Unlock()
		select {
		case <-ctx.Done():
			s.Lock()
			delete(a.notices, ch)
			s.putnotice(ch)
			if len(a.notices) == 0 && a.summary.CurVersion == 0 {
				delete(g, req.Appname)
			}
			if len(g) == 0 {
				delete(s.apps, req.Groupname)
			}
			s.Unlock()
			return nil, ctx.Err()
		case _, ok := <-ch:
			if !ok {
				return nil, cerror.ErrClosing
			}
		}
		s.Lock()
		delete(a.notices, ch)
		s.putnotice(ch)
		if int32(a.summary.CurVersion) != req.CurVersion {
			if len(a.notices) == 0 && a.summary.CurVersion == 0 {
				delete(g, req.Appname)
			}
			if len(g) == 0 {
				delete(s.apps, req.Groupname)
			}
			s.Unlock()
			return &api.WatchResp{
				Version:      int32(a.summary.CurVersion),
				AppConfig:    a.summary.CurAppConfig,
				SourceConfig: a.summary.CurSourceConfig,
			}, nil
		}
	}
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
