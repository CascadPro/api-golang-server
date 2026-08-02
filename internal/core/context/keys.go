package core_context

type CtxKey int

var (
	CtxKeyUserID    = CtxKey(0)
	CtxKeyUserRole  = CtxKey(1)
	CtxKeySessionID = CtxKey(2)
	CtxKeyMimeType  = CtxKey(3)
	CtxKeyTag       = CtxKey(4)
	CtxKeyIP        = CtxKey(5)
)
