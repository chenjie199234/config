package ecode

import (
	"net/http"

	cerror "github.com/chenjie199234/Corelib/error"
)

var (
	ErrUnknown  = cerror.ErrUnknown  //10000
	ErrReq      = cerror.ErrReq      //10001 // http code 400
	ErrResp     = cerror.ErrResp     //10002 // http code 500
	ErrSystem   = cerror.ErrSystem   //10003 // http code 500
	ErrAuth     = cerror.ErrAuth     //10004 // http code 401
	ErrLimit    = cerror.ErrLimit    //10005 // http code 503
	ErrBan      = cerror.ErrBan      //10006 // http code 403
	ErrNotExist = cerror.ErrNotExist //10007 // http code 404

	ErrBusiness1 = cerror.MakeError(20001, http.StatusBadRequest, "business err message")
)
