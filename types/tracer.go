package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// callLog is the result of LOG opCode
type CallLog struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    hexutil.Bytes  `json:"data"`
}

type CallTrace struct {
	From         common.Address  `json:"from"`
	Gas          *hexutil.Uint64 `json:"gas"`
	GasUsed      *hexutil.Uint64 `json:"gasUsed"`
	To           *common.Address `json:"to,omitempty"`
	Input        hexutil.Bytes   `json:"input"`
	Output       hexutil.Bytes   `json:"output,omitempty"`
	Error        string          `json:"error,omitempty"`
	RevertReason string          `json:"revertReason,omitempty"`
	Calls        []CallTrace     `json:"calls,omitempty"`
	Logs         []CallLog       `json:"logs,omitempty"`
	Value        *hexutil.Big    `json:"value,omitempty"`
	// Gencodec adds overridden fields at the end
	Type string `json:"type"`
}

type CallTraceResponse struct {
	Result CallTrace `json:"result"`
}

//func (c *CallTraceResponse) MarshalJSON() ([]byte, error) {
//	return json.Marshal(c)
//}
//
//func (c *CallTraceResponse) UnmarshalJSON(data []byte) error {
//	return json.Unmarshal(data, &c)
//}
