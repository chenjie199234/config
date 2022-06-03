// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.1
// 	protoc              v3.19.4
// source: api/admin.proto

package api

import (
	context "context"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	error1 "github.com/chenjie199234/Corelib/error"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
)

var _CGrpcPathAdminLogin = "/config.admin/login"
var _CGrpcPathAdminSearchUser = "/config.admin/search_user"
var _CGrpcPathAdminUpdateUserPermission = "/config.admin/update_user_permission"
var _CGrpcPathAdminAddNode = "/config.admin/add_node"
var _CGrpcPathAdminUpdateNode = "/config.admin/update_node"
var _CGrpcPathAdminMoveNode = "/config.admin/move_node"
var _CGrpcPathAdminDelNode = "/config.admin/del_node"
var _CGrpcPathAdminListNode = "/config.admin/list_node"
var _CGrpcPathAdminCheck = "/config.admin/check"

type AdminCGrpcClient interface {
	Login(context.Context, *LoginReq) (*LoginResp, error)
	SearchUser(context.Context, *SearchUserReq) (*SearchUserResp, error)
	UpdateUserPermission(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)
	AddNode(context.Context, *AddNodeReq) (*AddNodeResp, error)
	UpdateNode(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)
	MoveNode(context.Context, *MoveNodeReq) (*MoveNodeResp, error)
	DelNode(context.Context, *DelNodeReq) (*DelNodeResp, error)
	ListNode(context.Context, *ListNodeReq) (*ListNodeResp, error)
	Check(context.Context, *CheckReq) (*CheckResp, error)
}

type adminCGrpcClient struct {
	cc *cgrpc.CGrpcClient
}

func NewAdminCGrpcClient(c *cgrpc.CGrpcClient) AdminCGrpcClient {
	return &adminCGrpcClient{cc: c}
}

func (c *adminCGrpcClient) Login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(LoginResp)
	if e := c.cc.Call(ctx, _CGrpcPathAdminLogin, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *adminCGrpcClient) SearchUser(ctx context.Context, req *SearchUserReq) (*SearchUserResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(SearchUserResp)
	if e := c.cc.Call(ctx, _CGrpcPathAdminSearchUser, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *adminCGrpcClient) UpdateUserPermission(ctx context.Context, req *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(UpdateUserPermissionResp)
	if e := c.cc.Call(ctx, _CGrpcPathAdminUpdateUserPermission, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *adminCGrpcClient) AddNode(ctx context.Context, req *AddNodeReq) (*AddNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(AddNodeResp)
	if e := c.cc.Call(ctx, _CGrpcPathAdminAddNode, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *adminCGrpcClient) UpdateNode(ctx context.Context, req *UpdateNodeReq) (*UpdateNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(UpdateNodeResp)
	if e := c.cc.Call(ctx, _CGrpcPathAdminUpdateNode, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *adminCGrpcClient) MoveNode(ctx context.Context, req *MoveNodeReq) (*MoveNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(MoveNodeResp)
	if e := c.cc.Call(ctx, _CGrpcPathAdminMoveNode, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *adminCGrpcClient) DelNode(ctx context.Context, req *DelNodeReq) (*DelNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(DelNodeResp)
	if e := c.cc.Call(ctx, _CGrpcPathAdminDelNode, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *adminCGrpcClient) ListNode(ctx context.Context, req *ListNodeReq) (*ListNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(ListNodeResp)
	if e := c.cc.Call(ctx, _CGrpcPathAdminListNode, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *adminCGrpcClient) Check(ctx context.Context, req *CheckReq) (*CheckResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(CheckResp)
	if e := c.cc.Call(ctx, _CGrpcPathAdminCheck, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}

type AdminCGrpcServer interface {
	Login(context.Context, *LoginReq) (*LoginResp, error)
	SearchUser(context.Context, *SearchUserReq) (*SearchUserResp, error)
	UpdateUserPermission(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)
	AddNode(context.Context, *AddNodeReq) (*AddNodeResp, error)
	UpdateNode(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)
	MoveNode(context.Context, *MoveNodeReq) (*MoveNodeResp, error)
	DelNode(context.Context, *DelNodeReq) (*DelNodeResp, error)
	ListNode(context.Context, *ListNodeReq) (*ListNodeResp, error)
	Check(context.Context, *CheckReq) (*CheckResp, error)
}

func _Admin_Login_CGrpcHandler(handler func(context.Context, *LoginReq) (*LoginResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(LoginReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(LoginResp)
		}
		ctx.Write(resp)
	}
}
func _Admin_SearchUser_CGrpcHandler(handler func(context.Context, *SearchUserReq) (*SearchUserResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SearchUserReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SearchUserResp)
		}
		ctx.Write(resp)
	}
}
func _Admin_UpdateUserPermission_CGrpcHandler(handler func(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateUserPermissionReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/update_user_permission]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateUserPermissionResp)
		}
		ctx.Write(resp)
	}
}
func _Admin_AddNode_CGrpcHandler(handler func(context.Context, *AddNodeReq) (*AddNodeResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(AddNodeReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/add_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(AddNodeResp)
		}
		ctx.Write(resp)
	}
}
func _Admin_UpdateNode_CGrpcHandler(handler func(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateNodeReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/update_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateNodeResp)
		}
		ctx.Write(resp)
	}
}
func _Admin_MoveNode_CGrpcHandler(handler func(context.Context, *MoveNodeReq) (*MoveNodeResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(MoveNodeReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/move_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(MoveNodeResp)
		}
		ctx.Write(resp)
	}
}
func _Admin_DelNode_CGrpcHandler(handler func(context.Context, *DelNodeReq) (*DelNodeResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelNodeReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/del_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelNodeResp)
		}
		ctx.Write(resp)
	}
}
func _Admin_ListNode_CGrpcHandler(handler func(context.Context, *ListNodeReq) (*ListNodeResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(ListNodeReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/list_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(ListNodeResp)
		}
		ctx.Write(resp)
	}
}
func _Admin_Check_CGrpcHandler(handler func(context.Context, *CheckReq) (*CheckResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(CheckReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/check]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(CheckResp)
		}
		ctx.Write(resp)
	}
}
func RegisterAdminCGrpcServer(engine *cgrpc.CGrpcServer, svc AdminCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	//avoid lint
	_ = allmids
	engine.RegisterHandler("config.admin", "login", _Admin_Login_CGrpcHandler(svc.Login))
	engine.RegisterHandler("config.admin", "search_user", _Admin_SearchUser_CGrpcHandler(svc.SearchUser))
	engine.RegisterHandler("config.admin", "update_user_permission", _Admin_UpdateUserPermission_CGrpcHandler(svc.UpdateUserPermission))
	engine.RegisterHandler("config.admin", "add_node", _Admin_AddNode_CGrpcHandler(svc.AddNode))
	engine.RegisterHandler("config.admin", "update_node", _Admin_UpdateNode_CGrpcHandler(svc.UpdateNode))
	engine.RegisterHandler("config.admin", "move_node", _Admin_MoveNode_CGrpcHandler(svc.MoveNode))
	engine.RegisterHandler("config.admin", "del_node", _Admin_DelNode_CGrpcHandler(svc.DelNode))
	engine.RegisterHandler("config.admin", "list_node", _Admin_ListNode_CGrpcHandler(svc.ListNode))
	engine.RegisterHandler("config.admin", "check", _Admin_Check_CGrpcHandler(svc.Check))
}
