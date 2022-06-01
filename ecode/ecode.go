package ecode

import (
	"net/http"

	cerror "github.com/chenjie199234/Corelib/error"
)

var (
	ErrUnknown    = cerror.ErrUnknown    //10000
	ErrReq        = cerror.ErrReq        //10001 // http code 400
	ErrResp       = cerror.ErrResp       //10002 // http code 500
	ErrSystem     = cerror.ErrSystem     //10003 // http code 500
	ErrAuth       = cerror.ErrAuth       //10004 // http code 401
	ErrPermission = cerror.ErrPermission //10005 // http code 403
	ErrTooFast    = cerror.ErrTooFast    //10006 // http code 403
	ErrBan        = cerror.ErrBan        //10007 // http code 403
	ErrBusy       = cerror.ErrBusy       //10008 // http code 503
	ErrNotExist   = cerror.ErrNotExist   //10009 // http code 404

	ErrAppNotExist     = cerror.MakeError(20001, http.StatusBadRequest, "app doesn't exist")
	ErrAppAlreadyExist = cerror.MakeError(20002, http.StatusBadRequest, "app already exist")
	ErrWrongCipher     = cerror.MakeError(20003, http.StatusBadRequest, "wrong cipher")
	ErrCipherLength    = cerror.MakeError(20004, http.StatusBadRequest, "cipher must be empty or 32 byte length")
	ErrIndexNotExist   = cerror.MakeError(20005, http.StatusBadRequest, "config index doesn't exist")
	ErrConfigFormat    = cerror.MakeError(20001, http.StatusBadRequest, "config must use json object format")
)
