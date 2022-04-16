// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.1
// 	protoc             v3.19.4
// source: api/config.proto

package api

import (
	context "context"
	crpc "github.com/chenjie199234/Corelib/crpc"
	error1 "github.com/chenjie199234/Corelib/error"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	proto "google.golang.org/protobuf/proto"
)

var _CrpcPathConfigGroups = "/config.config/groups"
var _CrpcPathConfigApps = "/config.config/apps"
var _CrpcPathConfigCreate = "/config.config/create"
var _CrpcPathConfigUpdatecipher = "/config.config/updatecipher"
var _CrpcPathConfigGet = "/config.config/get"
var _CrpcPathConfigSet = "/config.config/set"
var _CrpcPathConfigRollback = "/config.config/rollback"
var _CrpcPathConfigWatch = "/config.config/watch"

type ConfigCrpcClient interface {
	//get all groups
	Groups(context.Context, *GroupsReq) (*GroupsResp, error)
	//get all apps in specific group
	Apps(context.Context, *AppsReq) (*AppsResp, error)
	//create one specific app
	Create(context.Context, *CreateReq) (*CreateResp, error)
	//update one specific app's cipher
	Updatecipher(context.Context, *UpdatecipherReq) (*UpdatecipherResp, error)
	//get one specific app's config
	Get(context.Context, *GetReq) (*GetResp, error)
	//set one specific app's config
	Set(context.Context, *SetReq) (*SetResp, error)
	//rollback one specific app's config
	Rollback(context.Context, *RollbackReq) (*RollbackResp, error)
	//watch on specific app's config
	Watch(context.Context, *WatchReq) (*WatchResp, error)
}

type configCrpcClient struct {
	cc *crpc.CrpcClient
}

func NewConfigCrpcClient(c *crpc.CrpcClient) ConfigCrpcClient {
	return &configCrpcClient{cc: c}
}

func (c *configCrpcClient) Groups(ctx context.Context, req *GroupsReq) (*GroupsResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigGroups, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(GroupsResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Apps(ctx context.Context, req *AppsReq) (*AppsResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigApps, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(AppsResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Create(ctx context.Context, req *CreateReq) (*CreateResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigCreate, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(CreateResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Updatecipher(ctx context.Context, req *UpdatecipherReq) (*UpdatecipherResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigUpdatecipher, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdatecipherResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Get(ctx context.Context, req *GetReq) (*GetResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigGet, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(GetResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Set(ctx context.Context, req *SetReq) (*SetResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigSet, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(SetResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Rollback(ctx context.Context, req *RollbackReq) (*RollbackResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigRollback, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(RollbackResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Watch(ctx context.Context, req *WatchReq) (*WatchResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigWatch, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(WatchResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}

type ConfigCrpcServer interface {
	//get all groups
	Groups(context.Context, *GroupsReq) (*GroupsResp, error)
	//get all apps in specific group
	Apps(context.Context, *AppsReq) (*AppsResp, error)
	//create one specific app
	Create(context.Context, *CreateReq) (*CreateResp, error)
	//update one specific app's cipher
	Updatecipher(context.Context, *UpdatecipherReq) (*UpdatecipherResp, error)
	//get one specific app's config
	Get(context.Context, *GetReq) (*GetResp, error)
	//set one specific app's config
	Set(context.Context, *SetReq) (*SetResp, error)
	//rollback one specific app's config
	Rollback(context.Context, *RollbackReq) (*RollbackResp, error)
	//watch on specific app's config
	Watch(context.Context, *WatchReq) (*WatchResp, error)
}

func _Config_Groups_CrpcHandler(handler func(context.Context, *GroupsReq) (*GroupsResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(GroupsReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GroupsResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Apps_CrpcHandler(handler func(context.Context, *AppsReq) (*AppsResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(AppsReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.config/apps]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(AppsResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Create_CrpcHandler(handler func(context.Context, *CreateReq) (*CreateResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(CreateReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.config/create]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(CreateResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Updatecipher_CrpcHandler(handler func(context.Context, *UpdatecipherReq) (*UpdatecipherResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(UpdatecipherReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.config/updatecipher]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdatecipherResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Get_CrpcHandler(handler func(context.Context, *GetReq) (*GetResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(GetReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.config/get]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Set_CrpcHandler(handler func(context.Context, *SetReq) (*SetResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(SetReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.config/set]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SetResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Rollback_CrpcHandler(handler func(context.Context, *RollbackReq) (*RollbackResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(RollbackReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.config/rollback]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(RollbackResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Watch_CrpcHandler(handler func(context.Context, *WatchReq) (*WatchResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(WatchReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.config/watch]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(WatchResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func RegisterConfigCrpcServer(engine *crpc.CrpcServer, svc ConfigCrpcServer, allmids map[string]crpc.OutsideHandler) {
	//avoid lint
	_ = allmids
	engine.RegisterHandler(_CrpcPathConfigGroups, _Config_Groups_CrpcHandler(svc.Groups))
	engine.RegisterHandler(_CrpcPathConfigApps, _Config_Apps_CrpcHandler(svc.Apps))
	engine.RegisterHandler(_CrpcPathConfigCreate, _Config_Create_CrpcHandler(svc.Create))
	engine.RegisterHandler(_CrpcPathConfigUpdatecipher, _Config_Updatecipher_CrpcHandler(svc.Updatecipher))
	engine.RegisterHandler(_CrpcPathConfigGet, _Config_Get_CrpcHandler(svc.Get))
	engine.RegisterHandler(_CrpcPathConfigSet, _Config_Set_CrpcHandler(svc.Set))
	engine.RegisterHandler(_CrpcPathConfigRollback, _Config_Rollback_CrpcHandler(svc.Rollback))
	engine.RegisterHandler(_CrpcPathConfigWatch, _Config_Watch_CrpcHandler(svc.Watch))
}
