// Code generated by protoc-gen-go-web. DO NOT EDIT.
// version:
// 	protoc-gen-go-web v0.0.1
// 	protoc            v3.19.4
// source: api/admin.proto

package api

import (
	context "context"
	error1 "github.com/chenjie199234/Corelib/error"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	pool "github.com/chenjie199234/Corelib/pool"
	web "github.com/chenjie199234/Corelib/web"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	http "net/http"
	strings "strings"
)

var _WebPathAdminLogin = "/config.admin/login"
var _WebPathAdminAddNode = "/config.admin/add_node"
var _WebPathAdminUpdateNode = "/config.admin/update_node"
var _WebPathAdminDelNode = "/config.admin/del_node"
var _WebPathAdminListNode = "/config.admin/list_node"

type AdminWebClient interface {
	Login(context.Context, *LoginReq, http.Header) (*LoginResp, error)
	AddNode(context.Context, *AddNodeReq, http.Header) (*AddNodeResp, error)
	UpdateNode(context.Context, *UpdateNodeReq, http.Header) (*UpdateNodeResp, error)
	DelNode(context.Context, *DelNodeReq, http.Header) (*DelNodeResp, error)
	ListNode(context.Context, *ListNodeReq, http.Header) (*ListNodeResp, error)
}

type adminWebClient struct {
	cc *web.WebClient
}

func NewAdminWebClient(c *web.WebClient) AdminWebClient {
	return &adminWebClient{cc: c}
}

func (c *adminWebClient) Login(ctx context.Context, req *LoginReq, header http.Header) (*LoginResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	data, e := c.cc.Post(ctx, _WebPathAdminLogin, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	resp := new(LoginResp)
	if len(data) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(data, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *adminWebClient) AddNode(ctx context.Context, req *AddNodeReq, header http.Header) (*AddNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	data, e := c.cc.Post(ctx, _WebPathAdminAddNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	resp := new(AddNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(data, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *adminWebClient) UpdateNode(ctx context.Context, req *UpdateNodeReq, header http.Header) (*UpdateNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	data, e := c.cc.Post(ctx, _WebPathAdminUpdateNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	resp := new(UpdateNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(data, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *adminWebClient) DelNode(ctx context.Context, req *DelNodeReq, header http.Header) (*DelNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	data, e := c.cc.Post(ctx, _WebPathAdminDelNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	resp := new(DelNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(data, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *adminWebClient) ListNode(ctx context.Context, req *ListNodeReq, header http.Header) (*ListNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	data, e := c.cc.Post(ctx, _WebPathAdminListNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	resp := new(ListNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(data, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}

type AdminWebServer interface {
	Login(context.Context, *LoginReq) (*LoginResp, error)
	AddNode(context.Context, *AddNodeReq) (*AddNodeResp, error)
	UpdateNode(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)
	DelNode(context.Context, *DelNodeReq) (*DelNodeResp, error)
	ListNode(context.Context, *ListNodeReq) (*ListNodeResp, error)
}

func _Admin_Login_WebHandler(handler func(context.Context, *LoginReq) (*LoginResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(LoginReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				ctx.Abort(error1.ErrReq)
				return
			}
			data := pool.GetBuffer()
			defer pool.PutBuffer(data)
			data.AppendByte('{')
			data.AppendByte('}')
			if data.Len() > 2 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data.Bytes(), req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		}
		resp, e := handler(ctx, req)
		ee := error1.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(LoginResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Admin_AddNode_WebHandler(handler func(context.Context, *AddNodeReq) (*AddNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(AddNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				ctx.Abort(error1.ErrReq)
				return
			}
			data := pool.GetBuffer()
			defer pool.PutBuffer(data)
			data.AppendByte('{')
			data.AppendString("\"pnode_id\":")
			if forms := ctx.GetForms("pnode_id"); len(forms) == 0 {
				data.AppendString("null")
			} else {
				data.AppendByte('[')
				for _, form := range forms {
					if len(form) == 0 {
						data.AppendString("0")
					} else {
						data.AppendString(form)
					}
					data.AppendByte(',')
				}
				data.Bytes()[data.Len()-1] = ']'
			}
			data.AppendByte(',')
			data.AppendString("\"node_name\":")
			if form := ctx.GetForm("node_name"); len(form) == 0 {
				data.AppendString("\"\"")
			} else if len(form) < 2 || form[0] != '"' || form[len(form)-1] != '"' {
				data.AppendByte('"')
				data.AppendString(form)
				data.AppendByte('"')
			} else {
				data.AppendString(form)
			}
			data.AppendByte(',')
			data.AppendString("\"node_data\":")
			if form := ctx.GetForm("node_data"); len(form) == 0 {
				data.AppendString("\"\"")
			} else if len(form) < 2 || form[0] != '"' || form[len(form)-1] != '"' {
				data.AppendByte('"')
				data.AppendString(form)
				data.AppendByte('"')
			} else {
				data.AppendString(form)
			}
			data.AppendByte('}')
			if data.Len() > 2 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data.Bytes(), req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/add_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := error1.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(AddNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Admin_UpdateNode_WebHandler(handler func(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				ctx.Abort(error1.ErrReq)
				return
			}
			data := pool.GetBuffer()
			defer pool.PutBuffer(data)
			data.AppendByte('{')
			data.AppendString("\"node_id\":")
			if forms := ctx.GetForms("node_id"); len(forms) == 0 {
				data.AppendString("null")
			} else {
				data.AppendByte('[')
				for _, form := range forms {
					if len(form) == 0 {
						data.AppendString("0")
					} else {
						data.AppendString(form)
					}
					data.AppendByte(',')
				}
				data.Bytes()[data.Len()-1] = ']'
			}
			data.AppendByte(',')
			data.AppendString("\"node_name\":")
			if form := ctx.GetForm("node_name"); len(form) == 0 {
				data.AppendString("\"\"")
			} else if len(form) < 2 || form[0] != '"' || form[len(form)-1] != '"' {
				data.AppendByte('"')
				data.AppendString(form)
				data.AppendByte('"')
			} else {
				data.AppendString(form)
			}
			data.AppendByte(',')
			data.AppendString("\"node_data\":")
			if form := ctx.GetForm("node_data"); len(form) == 0 {
				data.AppendString("\"\"")
			} else if len(form) < 2 || form[0] != '"' || form[len(form)-1] != '"' {
				data.AppendByte('"')
				data.AppendString(form)
				data.AppendByte('"')
			} else {
				data.AppendString(form)
			}
			data.AppendByte('}')
			if data.Len() > 2 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data.Bytes(), req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/update_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := error1.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(UpdateNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Admin_DelNode_WebHandler(handler func(context.Context, *DelNodeReq) (*DelNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(DelNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				ctx.Abort(error1.ErrReq)
				return
			}
			data := pool.GetBuffer()
			defer pool.PutBuffer(data)
			data.AppendByte('{')
			data.AppendString("\"node_id\":")
			if forms := ctx.GetForms("node_id"); len(forms) == 0 {
				data.AppendString("null")
			} else {
				data.AppendByte('[')
				for _, form := range forms {
					if len(form) == 0 {
						data.AppendString("0")
					} else {
						data.AppendString(form)
					}
					data.AppendByte(',')
				}
				data.Bytes()[data.Len()-1] = ']'
			}
			data.AppendByte('}')
			if data.Len() > 2 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data.Bytes(), req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/del_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := error1.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(DelNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Admin_ListNode_WebHandler(handler func(context.Context, *ListNodeReq) (*ListNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(ListNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				ctx.Abort(error1.ErrReq)
				return
			}
			data := pool.GetBuffer()
			defer pool.PutBuffer(data)
			data.AppendByte('{')
			data.AppendString("\"pnode_id\":")
			if forms := ctx.GetForms("pnode_id"); len(forms) == 0 {
				data.AppendString("null")
			} else {
				data.AppendByte('[')
				for _, form := range forms {
					if len(form) == 0 {
						data.AppendString("0")
					} else {
						data.AppendString(form)
					}
					data.AppendByte(',')
				}
				data.Bytes()[data.Len()-1] = ']'
			}
			data.AppendByte('}')
			if data.Len() > 2 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data.Bytes(), req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/config.admin/list_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := error1.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(ListNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func RegisterAdminWebServer(engine *web.WebServer, svc AdminWebServer, allmids map[string]web.OutsideHandler) {
	//avoid lint
	_ = allmids
	engine.Post(_WebPathAdminLogin, _Admin_Login_WebHandler(svc.Login))
	engine.Post(_WebPathAdminAddNode, _Admin_AddNode_WebHandler(svc.AddNode))
	engine.Post(_WebPathAdminUpdateNode, _Admin_UpdateNode_WebHandler(svc.UpdateNode))
	engine.Post(_WebPathAdminDelNode, _Admin_DelNode_WebHandler(svc.DelNode))
	engine.Post(_WebPathAdminListNode, _Admin_ListNode_WebHandler(svc.ListNode))
}
