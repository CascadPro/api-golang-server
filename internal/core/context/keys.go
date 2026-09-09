package core_context

type CtxKey string

var (
	ctxKeyUserID    = CtxKey("user_id")
	ctxKeyUserRole  = CtxKey("user_role")
	ctxKeySessionID = CtxKey("session_id")
	ctxKeyMimeType  = CtxKey("mime_type")
	ctxKeyTag       = CtxKey("tag")
	ctxKeyIP        = CtxKey("ip")
	ctxKeyRequestID = CtxKey("request_id")
	ctxKeyLocale    = CtxKey("locale")
)
