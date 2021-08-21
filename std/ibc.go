package std

type IBCMsg struct {
	Transfer     *TransferMsg     `json:"transfer,omitempty"`
	SendPacket   *SendPacketMsg   `json:"send_packet,omitempty"`
	CloseChannel *CloseChannelMsg `json:"close_channel,omitempty"`
}

type TransferMsg struct {
	ChannelID string `json:"channel_id"`
	ToAddress string `json:"to_address"`
	Amount    Coin   `json:"amount"`
	// Timeout   IBCTimeout `json:"timeout"`
}

type SendPacketMsg struct {
	ChannelID string `json:"channel_id"`
	Data      []byte `json:"data"`
	// Timeout   IBCTimeout `json:"timeout"`
}

type CloseChannelMsg struct {
	ChannelID string `json:"channel_id"`
}
