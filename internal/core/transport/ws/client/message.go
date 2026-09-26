package core_ws_client

type ChannelMessage struct {
	data        []byte
	closeAfter  bool
	closeCode   CloseCode
	closeReason CloseReason
}
