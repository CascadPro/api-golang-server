package core_context

type CtxKey string

var (
	CtxKeyUserID    = CtxKey("user_id")
	CtxKeyUserRole  = CtxKey("user_role")
	CtxKeySessionID = CtxKey("session_id")
	CtxKeyMimeType  = CtxKey("mime_type")
	CtxKeyTag       = CtxKey("tag")
	CtxKeyIP        = CtxKey("ip")
	CtxKeyRequestID = CtxKey("request_id")
)
