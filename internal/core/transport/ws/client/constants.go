package core_ws_client

type CloseReason string

const (
	CloseReasonSessionRevoked = CloseReason("session_revoked")
)

type CloseCode int

const (
	CloseCodeSessionRevoked = CloseCode(4001)
)

const (
	MaxMessageSize = 64 * 1024
)
