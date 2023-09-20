package types

import (
	"encoding/json"

	utilbellatrix "github.com/attestantio/go-eth2-client/util/bellatrix"
)

type TobTxsSubmitRequest struct {
	TobTxs     utilbellatrix.ExecutionPayloadTransactions
	Slot       uint64
	ParentHash string
}

func (t *TobTxsSubmitRequest) MarshalJSON() ([]byte, error) {
	txBytes, err := t.TobTxs.MarshalSSZ()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		TobTxs     []byte `json:"tobTxs"`
		Slot       uint64 `json:"slot"`
		ParentHash string `json:"parentHash"`
	}{
		TobTxs:     txBytes,
		Slot:       t.Slot,
		ParentHash: t.ParentHash,
	})
}

func (t *TobTxsSubmitRequest) UnmarshalJSON(data []byte) error {
	var intermediateJson struct {
		TobTxs     []byte `json:"tobTxs"`
		Slot       uint64 `json:"slot"`
		ParentHash string `json:"parentHash"`
	}
	err := json.Unmarshal(data, &intermediateJson)
	if err != nil {
		return err
	}

	err = t.TobTxs.UnmarshalSSZ(intermediateJson.TobTxs)
	if err != nil {
		return err
	}
	t.Slot = intermediateJson.Slot
	t.ParentHash = intermediateJson.ParentHash

	return nil
}
