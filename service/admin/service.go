package admin

import (
	"context"
	"time"

	"github.com/chenjie199234/config/api"
	"github.com/chenjie199234/config/config"
	admindao "github.com/chenjie199234/config/dao/admin"
	"github.com/chenjie199234/config/ecode"

	cerror "github.com/chenjie199234/Corelib/error"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/log"
	//"github.com/chenjie199234/Corelib/web"
)

//Service subservice for admin business
type Service struct {
	adminDao *admindao.Dao
}

//Start -
func Start() *Service {
	return &Service{
		adminDao: admindao.NewDao(config.GetSql("config_sql"), config.GetRedis("config_redis"), config.GetMongo("config_mongo")),
	}
}

func (s *Service) Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error) {
	var userid string
	//TODO get userid
	start := time.Now()
	end := start.Add(config.AC.TokenExpire.StdDuration())
	tokenstr := publicmids.MakeToken(config.AC.TokenSecret, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, userid, uint64(start.Unix()), uint64(end.Unix()))
	return &api.LoginResp{Token: tokenstr}, nil
}
func (s *Service) SearchUser(ctx context.Context, req *api.SearchUserReq) (*api.SearchUserResp, error) {
	//TODO
	return &api.SearchUserResp{}, nil
}
func (s *Service) UpdateUserPermission(ctx context.Context, req *api.UpdateUserPermissionReq) (*api.UpdateUserPermissionResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	if !req.Admin && req.Canwrite && !req.Canread {
		return nil, ecode.ErrReq
	}
	targetobj, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[UpdateUserPermission] target userid:", req.UserId, "format error:", e)
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[UpdateUserPermission] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.adminDao.MongoUpdateUserPermission(ctx, obj, targetobj, req.NodeId, req.Admin, req.Canread, req.Canwrite); e != nil {
		log.Error(ctx, "[UpdateUserPermission] userid:", userid, "target userid:", req.UserId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.UpdateUserPermissionResp{}, nil
}

func (s *Service) AddNode(ctx context.Context, req *api.AddNodeReq) (*api.AddNodeResp, error) {
	if req.PnodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[AddNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.adminDao.MongoAddNode(ctx, obj, req.PnodeId, req.NodeName, req.NodeData); e != nil {
		log.Error(ctx, "[AddNode] userid:", userid, "name:", req.NodeName, "data:", req.NodeData, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.AddNodeResp{}, nil
}
func (s *Service) UpdateNode(ctx context.Context, req *api.UpdateNodeReq) (*api.UpdateNodeResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[UpdateNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.adminDao.MongoUpdateNode(ctx, obj, req.NodeId, req.NodeName, req.NodeData); e != nil {
		log.Error(ctx, "[UpdateNode] userid:", userid, "nodeid:", req.NodeId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.UpdateNodeResp{}, nil
}
func (s *Service) MoveNode(ctx context.Context, req *api.MoveNodeReq) (*api.MoveNodeResp, error) {
	if req.NodeId[0] != 0 || req.PnodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[MoveNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	if e := s.adminDao.MongoMoveNode(ctx, obj, req.NodeId, req.PnodeId); e != nil {
		log.Error(ctx, "[MoveNode] userid:", userid, "old nodeid:", req.NodeId, "new parent:", req.PnodeId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.MoveNodeResp{}, nil
}
func (s *Service) DelNode(ctx context.Context, req *api.DelNodeReq) (*api.DelNodeResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[DelNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.adminDao.MongoDeleteNode(ctx, obj, req.NodeId); e != nil {
		log.Error(ctx, "[DelNode] userid:", userid, "nodeid:", req.NodeId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.DelNodeResp{}, nil
}
func (s *Service) ListNode(ctx context.Context, req *api.ListNodeReq) (*api.ListNodeResp, error) {
	if req.PnodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[ListNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	nodes, e := s.adminDao.MongoListNode(ctx, obj, req.PnodeId)
	if e != nil {
		log.Error(ctx, "[ListNode] userid:", userid, "pnodeid:", req.PnodeId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	resp := &api.ListNodeResp{
		Nodes: make([]*api.NodeInfo, 0, len(nodes)),
	}
	for _, node := range nodes {
		resp.Nodes = append(resp.Nodes, &api.NodeInfo{
			NodeId:   node.NodeId,
			NodeName: node.NodeName,
			NodeData: node.NodeData,
		})
	}
	return resp, nil
}
func (s *Service) Check(ctx context.Context, req *api.CheckReq) (*api.CheckResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	obj, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[Check] userid:", req.UserId, "format error:", e)
		return nil, ecode.ErrReq
	}
	canread, canwrite, admin, e := s.adminDao.MongoGetUserPermission(ctx, obj, req.NodeId)
	if e != nil {
		log.Error(ctx, "[Check] userid:", req.UserId, "nodeid:", req.NodeId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.CheckResp{Canread: canread, Canwrite: canwrite, Admin: admin}, nil
}

//Stop -
func (s *Service) Stop() {

}
